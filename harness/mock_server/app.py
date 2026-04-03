"""FastAPI application for the mock Anthropic API server."""

import json
import uuid
from typing import Any

from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse

from .scenarios import registry, Scenario
from .recorder import recorder


app = FastAPI(title="Mock Anthropic API")


def generate_message_id() -> str:
    """Generate a unique message ID."""
    return f"msg_{uuid.uuid4().hex[:12]}"


def generate_tool_use_id() -> str:
    """Generate a unique tool_use ID."""
    return f"toolu_{uuid.uuid4().hex[:10]}"


def get_current_response(scenario: Scenario) -> dict[str, Any] | None:
    """Get the current response based on turn count."""
    turn_count = registry.get_turn_count(scenario.name)
    if turn_count > 0:
        turn_index = (turn_count - 1) % len(scenario.messages)
    else:
        turn_index = 0
    return scenario.messages[turn_index]


def build_message_start_event(message_id: str) -> str:
    """Build the message_start SSE event."""
    data = {
        "type": "message_start",
        "message": {
            "id": message_id,
            "type": "message",
            "role": "assistant",
            "content": [],
            "model": "claude-sonnet-4-20250514",
            "stop_reason": None,
            "stop_sequence": None,
            "usage": {
                "input_tokens": 10,
                "output_tokens": 0
            }
        }
    }
    return f"event: message_start\ndata: {json.dumps(data)}\n\n"


def build_content_block_start_event(index: int, block_type: str, name: str = None, tool_id: str = None) -> str:
    """Build the content_block_start SSE event."""
    if block_type == "text":
        content_block = {"type": "text", "text": ""}
    elif block_type == "tool_use":
        content_block = {
            "type": "tool_use",
            "id": tool_id,
            "name": name,
            "input": {}
        }
    else:
        content_block = {"type": block_type}

    data = {
        "type": "content_block_start",
        "index": index,
        "content_block": content_block
    }
    return f"event: content_block_start\ndata: {json.dumps(data)}\n\n"


def build_text_delta_event(index: int, text: str) -> str:
    """Build content_block_delta with text_delta."""
    data = {
        "type": "content_block_delta",
        "index": index,
        "delta": {
            "type": "text_delta",
            "text": text
        }
    }
    return f"event: content_block_delta\ndata: {json.dumps(data)}\n\n"


def build_input_json_delta_event(index: int, partial_json: str) -> str:
    """Build content_block_delta with input_json_delta for tool_use."""
    data = {
        "type": "content_block_delta",
        "index": index,
        "delta": {
            "type": "input_json_delta",
            "partial_json": partial_json
        }
    }
    return f"event: content_block_delta\ndata: {json.dumps(data)}\n\n"


def build_content_block_stop_event(index: int) -> str:
    """Build the content_block_stop SSE event."""
    data = {
        "type": "content_block_stop",
        "index": index
    }
    return f"event: content_block_stop\ndata: {json.dumps(data)}\n\n"


def build_message_delta_event(stop_reason: str = "end_turn", output_tokens: int = 20) -> str:
    """Build the message_delta SSE event."""
    data = {
        "type": "message_delta",
        "delta": {
            "stop_reason": stop_reason,
            "stop_sequence": None
        },
        "usage": {
            "output_tokens": output_tokens
        }
    }
    return f"event: message_delta\ndata: {json.dumps(data)}\n\n"


def build_message_stop_event() -> str:
    """Build the message_stop SSE event."""
    return "event: message_stop\ndata: {}\n\n"


def generate_sse_events(response: dict[str, Any], message_id: str) -> list[str]:
    """Generate SSE events for a response block."""
    events = []

    # message_start
    events.append(build_message_start_event(message_id))

    # content_block_start
    block_type = response.get("type", "text")
    if block_type == "tool_use":
        tool_id = response.get("id", generate_tool_use_id())
        tool_name = response.get("name", "Unknown")
        events.append(build_content_block_start_event(0, "tool_use", tool_name, tool_id))
    else:
        events.append(build_content_block_start_event(0, "text"))

    # content_block_delta - stream text or tool input
    if block_type == "tool_use":
        # For tool_use, send the input as partial_json
        tool_input = response.get("input", {})
        partial_json = json.dumps(tool_input)
        events.append(build_input_json_delta_event(0, partial_json))
    else:
        # For text, stream the text content
        text = response.get("text", "")
        events.append(build_text_delta_event(0, text))

    # content_block_stop
    events.append(build_content_block_stop_event(0))

    # message_delta
    events.append(build_message_delta_event())

    # message_stop
    events.append(build_message_stop_event())

    return events


async def generate_stream(scenario_name: str, response: dict[str, Any]):
    """Generate the SSE stream."""
    message_id = generate_message_id()
    events = generate_sse_events(response, message_id)

    for event in events:
        yield event


@app.post("/v1/messages")
async def create_message(request: Request):
    """Handle POST /v1/messages - Anthropic Messages API."""
    # Parse request body
    body = await request.json()

    # Record the request
    recorder.record(body)

    # Get scenario from X-Scenario header or default to streaming_text
    scenario_name = request.headers.get("X-Scenario", "streaming_text")
    scenario = registry.get_scenario(scenario_name)

    if not scenario:
        scenario = registry.get_scenario("streaming_text")

    # Record this request for turn tracking
    registry.record_request(scenario_name, body)

    # Get current response
    response = get_current_response(scenario)
    if not response:
        response = {"type": "text", "text": "No more responses"}

    # Check if streaming
    stream = body.get("stream", False)

    if not stream:
        # Non-streaming response - generate all events at once
        message_id = generate_message_id()
        events = generate_sse_events(response, message_id)

        # Combine all events into a single response
        full_content = ""
        for event in events:
            # Parse the event to get the data portion
            if event.startswith("event: content_block_delta"):
                # Extract text from text_delta events
                lines = event.split("\n")
                for line in lines:
                    if line.startswith("data:"):
                        data = json.loads(line[5:])
                        if data.get("type") == "content_block_delta":
                            delta = data.get("delta", {})
                            if delta.get("type") == "text_delta":
                                full_content += delta.get("text", "")

        return {
            "id": message_id,
            "type": "message",
            "role": "assistant",
            "content": [{"type": "text", "text": full_content}],
            "model": body.get("model", "claude-sonnet-4-20250514"),
            "stop_reason": "end_turn",
            "stop_sequence": None,
            "usage": {
                "input_tokens": 10,
                "output_tokens": len(full_content) // 4
            }
        }
    else:
        # Streaming response
        return StreamingResponse(
            generate_stream(scenario_name, response),
            media_type="text/event-stream"
        )


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "ok"}


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": "Mock Anthropic API",
        "endpoints": ["/v1/messages", "/health"]
    }