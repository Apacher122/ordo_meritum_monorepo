#!/bin/bash

# Try to run Node first
ts-node -r tsconfig-paths/register src/index.ts

# If Node exits/crashes, fall back to a shell
echo "Node exited. Container is still alive."
exec bash
