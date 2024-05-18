#!/bin/bash
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -p|--patched-go-runtime)
      pathToPatchedGoRuntime="$2"
      shift
      shift
      ;;
    -g|--go-root)
      pathToGoRoot="$2"
      shift
      shift
      ;;
    -i|--overhead-inserter)
      pathToOverHeaderInserter="$2"
      shift
      shift
      ;;
    -r|--overhead-remover)
      pathToOverheadRemover="$2"
      shift
      shift
      ;;
    -f|--file)
      file="$2"
      shift
      shift
      ;;
    *)
      shift
      ;;
  esac
done

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverHeaderInserter" ] || [ -z "$pathToOverheadRemover" ] || [ -z "$file" ]; then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -f <file>"
  exit 1
fi

echo "Running full workflow on $file"

echo "Step 0: Remove Overhead just in case"
$pathToOverheadRemover -f $file
if [ $? -ne 0 ]; then
  echo "Error: Failed to remove overhead"
  exit 1
fi

echo "Step 1: Add Overhead to file"
$pathToOverHeaderInserter -f $file
if [ $? -ne 0 ]; then
  echo "Error: Failed to add overhead"
  exit 1
fi

echo "Step 2: Run with patched go runtime"
echo "Step 2.1: save current go root and set adjusted goroot"
export GOROOT=$pathToGoRoot
echo "Step 2.2: run program"
$pathToPatchedGoRuntime run $file
# Step 3: Analyze Trace
echo "Step 3: Analyze Trace"
echo "Step 3.1: Unset goroot"
unset GOROOT
echo "Step 3.3: Remove Overhead"
$pathToOverheadRemover -f $file
if [ $? -ne 0 ]; then
  echo "Error: Failed to remove overhead"
  exit 1
fi

echo "Workflow completed successfully."
exit 0
