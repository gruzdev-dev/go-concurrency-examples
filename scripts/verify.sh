#!/bin/bash

CONFIG_FILE="${1:-configs/race.json}"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

MODE=$(grep -o '"mode"[[:space:]]*:[[:space:]]*"[^"]*"' "$CONFIG_FILE" | sed 's/.*"\([^"]*\)"$/\1/')

case "$MODE" in
    "mutex_copy")
        BUILD_TAGS="-tags=mutex_copy"
        ;;
    "interface_tearing")
        BUILD_TAGS="-tags=interface_tearing"
        ;;
    *)
        BUILD_TAGS=""
        ;;
esac

echo "Config: $CONFIG_FILE"
echo "Mode: $MODE"
[ -n "$BUILD_TAGS" ] && echo "Tags: $BUILD_TAGS"
echo ""

echo "[1/2] Static analysis (go vet)"

set +e
if [ -n "$BUILD_TAGS" ]; then
    VET_OUTPUT=$(go vet $BUILD_TAGS ./... 2>&1)
else
    VET_OUTPUT=$(go vet ./... 2>&1)
fi
VET_EXIT=$?
set -e

if [ $VET_EXIT -ne 0 ] || [ -n "$VET_OUTPUT" ]; then
    echo -e "Result: ${GREEN}DETECTED${NC}"
    echo "$VET_OUTPUT"
    STATIC=1
else
    echo -e "Result: ${RED}NOT DETECTED${NC}"
    STATIC=0
fi

echo ""
echo "[2/2] Dynamic analysis (go run -race)"

set +e
if [ -n "$BUILD_TAGS" ]; then
    RACE_OUTPUT=$(go run -race $BUILD_TAGS . -config "$CONFIG_FILE" 2>&1)
else
    RACE_OUTPUT=$(go run -race . -config "$CONFIG_FILE" 2>&1)
fi
set -e

echo "$RACE_OUTPUT"
echo ""

if echo "$RACE_OUTPUT" | grep -q "WARNING: DATA RACE"; then
    echo -e "Result: ${GREEN}DETECTED${NC}"
    DYNAMIC=1
else
    echo -e "Result: ${RED}NOT DETECTED${NC}"
    DYNAMIC=0
fi

echo ""
echo "Summary"
echo -n "  go vet: "
[ $STATIC -eq 1 ] && echo -e "${GREEN}DETECTED${NC}" || echo -e "${RED}NOT DETECTED${NC}"
echo -n "  -race:  "
[ $DYNAMIC -eq 1 ] && echo -e "${GREEN}DETECTED${NC}" || echo -e "${RED}NOT DETECTED${NC}"

exit 0
