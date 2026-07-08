# Unit Test Observations — code unreachable or untestable by unit tests

This file tracks implementation code that cannot be (fully) unit-tested
through public entry points, per the rule: never add helpers/seams to
implementation code just to make it testable. Each entry needs manual
investigation by the maintainer.

Format: `path` — reason — suggested action.

## Untestable / unreachable code

(populated during the unit test refactoring, phases 1-8)

## Gio-UI-heavy files (covered by integration suite, not unit-testable)

(populated in phase 8)

## Dead code found while testing

(populated during the unit test refactoring)
