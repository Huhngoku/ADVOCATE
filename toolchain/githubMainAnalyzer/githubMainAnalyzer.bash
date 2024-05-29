#!/bin/bash
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -a|--advocate)
      pathToAdvocate="$2"
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

#intialize variables
pathToPatchedGoRuntime="$pathToAdvocate/go-patch/bin/go"
pathToGoRoot="$pathToAdvocate/go-patch"
pathToOverHeaderInserter="$pathToAdvocate/toolchain/overHeadInserter/inserter"
pathToOverheadRemover="$pathToAdvocate/toolchain/overHeadRemover/remover"
analyzer="$pathToAdvocate/analyzer/analyzer"

#print usage on misuse
if [ -z "$pathToAdvocate" ]; then
    echo "Path to advocate is empty"
    exit 1
fi
if [ -z "$githubUrl" ]; then
    echo "Github URL is empty"
    exit 1
fi
echo "cloning repository: $githubUrl"
if ! git clone "$githubUrl"; then
    echo "error: failed to clone the repository."
    exit 1
fi
cd $(basename "$githubUrl" .git)

fileToExecute=$(find . -name "*.go" -exec grep -q "func main()" {} \; -print -quit)
if [ $? -ne 0 ]; then
    echo "Error: Failed to find the main Go file."
    exit 1
fi

fileToExecute=$(echo "$fileToExecute" | sed 's|^\./||')
#remove overhead just in case
echo "Removing overhead"
echo "$pathToOverheadRemover -f $fileToExecute"
"$pathToOverheadRemover" -f "$fileToExecute"
#insert overhead
echo "Inserting overhead"
echo "$pathToOverHeaderInserter -f $fileToExecute"
"$pathToOverHeaderInserter" -f "$fileToExecute"
#run main
echo "Running main"
echo "$pathToPatchedGoRuntime run $fileToExecute"
"$pathToPatchedGoRuntime" run "$fileToExecute"
#remove overhead
echo "Removing overhead"
echo "$pathToOverheadRemover -f $fileToExecute"
"$pathToOverheadRemover" -f "$fileToExecute"
echo "Run Analysis"
"$analyzer" -t advocateTrace
if [ $? -ne 0 ]; then
    echo "Error: Failed to run the analysis."
    exit 1
fi