"""Release hygiene and documentation truth checks."""

import subprocess
from pathlib import Path


ROOT = Path(__file__).parent.parent


def test_release_hygiene_docs_exist_and_define_evidence_contracts():
    hygiene = ROOT / "docs" / "release-hygiene.md"
    inventory = ROOT / "docs" / "docs-inventory.md"

    hygiene_text = hygiene.read_text()
    inventory_text = inventory.read_text()

    for expected in [
        "OpenSpec Hygiene Checklist",
        "Task Evidence Format",
        "Archive Readiness",
        "Release State Matrix",
        "Known Gaps",
    ]:
        assert expected in hygiene_text

    assert "Source of truth" in inventory_text
    assert "Generated output" in inventory_text
    assert "docs/.vitepress/dist" in inventory_text


def test_release_hygiene_script_checks_docs_and_binary():
    result = subprocess.run(
        ["./scripts/check-release-hygiene.sh"],
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        check=False,
    )

    assert result.returncode == 0, result.stdout
    assert "OpenSpec hygiene" in result.stdout
    assert "install smoke" in result.stdout
