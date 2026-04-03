"""Output quality evaluators for testing Go CLI responses."""

import re
from typing import List, Dict, Any


class TextCompletenessEvaluator:
    """Evaluates whether actual text contains expected text content.
    
    Can check for exact substring match or ratio of expected words found.
    """
    
    def __init__(self, expected_text: str):
        """Initialize with expected text content.
        
        Args:
            expected_text: The text that should be present in the actual output.
        """
        self.expected_text = expected_text
        self.expected_words = set(expected_text.lower().split())
    
    def evaluate(self, actual_text: str) -> Dict[str, Any]:
        """Evaluate if actual text contains expected text.
        
        Args:
            actual_text: The actual output text to evaluate.
            
        Returns:
            Dictionary with:
                - complete: bool - whether all expected content is present
                - missing_ratio: float - ratio of missing words (0.0 to 1.0)
                - missing_parts: list - list of missing words/phrases
        """
        if not self.expected_text:
            return {
                "complete": True,
                "missing_ratio": 0.0,
                "missing_parts": []
            }
        
        if not actual_text:
            return {
                "complete": False,
                "missing_ratio": 1.0,
                "missing_parts": [self.expected_text]
            }
        
        # Check for exact substring match first
        if self.expected_text in actual_text:
            return {
                "complete": True,
                "missing_ratio": 0.0,
                "missing_parts": []
            }
        
        # Check word-level match ratio
        actual_words = set(actual_text.lower().split())
        found_words = self.expected_words & actual_words
        
        if len(self.expected_words) == 0:
            word_ratio = 1.0
        else:
            word_ratio = len(found_words) / len(self.expected_words)
        
        missing_words = self.expected_words - found_words
        missing_parts = list(missing_words)
        
        # Also check for longer phrases that might be missing
        phrase_length = 3
        expected_phrase_words = self.expected_text.lower().split()
        for i in range(len(expected_phrase_words) - phrase_length + 1):
            phrase = " ".join(expected_phrase_words[i:i + phrase_length])
            if phrase not in actual_text.lower():
                if phrase not in missing_parts:
                    missing_parts.append(phrase)
        
        complete = word_ratio >= 0.9
        
        return {
            "complete": complete,
            "missing_ratio": 1.0 - word_ratio,
            "missing_parts": missing_parts[:10]  # Limit to first 10 missing parts
        }


class ToolCallCorrectnessEvaluator:
    """Evaluates whether actual tool calls match expected tool calls.
    
    Compares both tool names and parameters.
    """
    
    def __init__(self, expected_calls: List[Dict[str, Any]]):
        """Initialize with expected tool calls.
        
        Args:
            expected_calls: List of expected tool calls, each with 'name' and 'params'.
        """
        self.expected_calls = expected_calls
    
    def _normalize_params(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Normalize parameters for comparison."""
        if not params:
            return {}
        # Convert to dict if it's a different type
        return dict(params)
    
    def _params_equal(self, expected: Dict[str, Any], actual: Dict[str, Any]) -> bool:
        """Check if parameters are equal."""
        expected_norm = self._normalize_params(expected)
        actual_norm = self._normalize_params(actual)
        
        if expected_norm == actual_norm:
            return True
        
        # Check if actual contains all expected params
        for key, value in expected_norm.items():
            if key not in actual_norm:
                return False
            actual_value = actual_norm[key]
            # Handle string containment
            if isinstance(value, str) and isinstance(actual_value, str):
                if value not in actual_value:
                    return False
            elif value != actual_value:
                return False
        
        return True
    
    def _find_matching_call(self, expected_call: Dict[str, Any], 
                           actual_calls: List[Dict[str, Any]]) -> tuple[bool, int]:
        """Find a matching actual call for the expected call.
        
        Returns:
            Tuple of (found, index) where index is -1 if not found.
        """
        expected_name = expected_call.get("name", "")
        
        for i, actual in enumerate(actual_calls):
            actual_name = actual.get("name", "")
            if expected_name == actual_name:
                # Check params
                expected_params = expected_call.get("params", {})
                actual_params = actual.get("params", {})
                
                if self._params_equal(expected_params, actual_params):
                    return True, i
        
        return False, -1
    
    def evaluate(self, actual_calls: List[Dict[str, Any]]) -> Dict[str, Any]:
        """Evaluate if actual tool calls match expected.
        
        Args:
            actual_calls: List of actual tool calls to compare.
            
        Returns:
            Dictionary with:
                - correct: bool - whether all expected calls match actual calls
                - mismatches: list - list of mismatch descriptions
        """
        mismatches = []
        
        # Check for missing or extra calls
        remaining_actual = list(actual_calls)
        
        for i, expected in enumerate(self.expected_calls):
            found, idx = self._find_matching_call(expected, remaining_actual)
            
            if found:
                # Remove from remaining to mark as matched
                remaining_actual.pop(idx)
            else:
                expected_name = expected.get("name", "unknown")
                expected_params = expected.get("params", {})
                mismatches.append({
                    "type": "missing",
                    "expected_index": i,
                    "expected_name": expected_name,
                    "expected_params": expected_params
                })
        
        # Check for extra calls
        for extra in remaining_actual:
            mismatches.append({
                "type": "extra",
                "actual_name": extra.get("name", "unknown"),
                "actual_params": extra.get("params", {})
            })
        
        correct = len(mismatches) == 0 and len(actual_calls) == len(self.expected_calls)
        
        return {
            "correct": correct,
            "mismatches": mismatches
        }