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
    -t|--gen-trace)
      genTrace="$2"
      shift 
      shift 
      ;;
    -a|--analyzer)
      analyzer="$2"
      shift 
      shift 
      ;;
    -u|--github-url)
      githubUrl="$2"
      shift 
      shift 
      ;;
    *)
      shift
      ;;
  esac
done

if [ -z "$genTrace" ] || [ -z "$analyzer" ] || [ -z "$githubUrl" ] || [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverHeaderInserter" ] || [ -z "$pathToOverheadRemover" ] ;then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -t <gen-trace> -a <analyzer> -u <github-url>"
  exit 1
fi

echo "Cloning repository: $githubUrl"
if ! git clone "$githubUrl"; then
    echo "Error: Failed to clone the repository."
    exit 1
fi
cd $(basename "$githubUrl" .git)

fileToExecute=$(find . -name "*.go" -exec grep -q "func main()" {} \; -print -quit)
if [ $? -ne 0 ]; then
    echo "Error: Failed to find the main Go file."
    exit 1
fi

fileToExecute=$(echo "$fileToExecute" | sed 's|^\./||')
echo "Generating trace"
"$genTrace" -f "$fileToExecute" -p "$pathToPatchedGoRuntime" -g "$pathToGoRoot" -i "$pathToOverHeaderInserter" -r "$pathToOverheadRemover"
if [ $? -ne 0 ]; then
    echo "Error: Failed to generate the trace."
    exit 1
fi
echo "Run Analysis"
"$analyzer" -t advocateTrace
if [ $? -ne 0 ]; then
    echo "Error: Failed to run the analysis."
    exit 1
fi