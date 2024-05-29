#!/bin/bash
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -a|--advocate)
      pathToAdvocate="$2"
      shift
      shift
      ;;
    -f|--folder)
      dir="$2"
      shift
      shift
      ;;
    *)
      shift
      ;;
  esac
done

if [ -z "$pathToAdvocate" ]; then
  echo "Path to advocate is empty"
  exit 1
fi
if [ -z "$dir" ]; then
  echo "Directory is empty"
  exit 1
fi

#Initialize Variables
pathToPatchedGoRuntime="$pathToAdvocate/go-patch/bin/go"
pathToGoRoot="$pathToAdvocate/go-patch"
pathToOverheadInserter="$pathToAdvocate/toolchain/unitTestOverheadInserter/unitTestOverheadInserter"
pathToOverheadRemover="$pathToAdvocate/toolchain/unitTestOverheadRemover/unitTestOverheadRemover"
pathToAnalyzer="$pathToAdvocate/analyzer/analyzer"
pathToFullWorkflowExecutor="$pathToAdvocate/toolchain/unitTestFullWorkflow/unitTestFullWorkflow.bash"

cd "$dir"
echo  "In directory: $dir"
mkdir -p advocation_results

test_files=$(find "$dir" -name "*_test.go")
total_files=$(echo "$test_files" | wc -l)
current_file=1
#echo "Test files: $test_files"
for file in $test_files; do
    echo "Progress: $current_file/$total_files"
    echo "Processing file: $file"
    package_path=$(dirname "$file")
    test_functions=$(grep -oE "[a-zA-Z0-9_]+ *Test[a-zA-Z0-9_]*" $file | sed 's/ *\(t *\*testing\.T\)//' | sed 's/func //')
    for test_func in $test_functions; do
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        #runfullworkflow for single test pass all 
        echo "Running full workflow for test: $test_func in package: $package_path in file: $file"
        adjustedPackagePath=$(echo "$package_path" | sed "s|$dir||g")
        #make folder for each test function output
        mkdir -p advocation_results/$packageName-$fileName-$test_func
        $pathToFullWorkflowExecutor -a $pathToAdvocate -p $adjustedPackagePath -f $dir -tf $file -t $test_func &> advocation_results/$packageName-$fileName-$test_func/advocation_results.txt
    done
    current_file=$((current_file+1))
done
