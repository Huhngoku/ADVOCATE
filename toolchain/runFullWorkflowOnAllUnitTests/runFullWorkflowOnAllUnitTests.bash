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
mkdir -p advocateResult

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
        echo "Running full workflow for test: $test_func in package: $package_path in file: $file"
        adjustedPackagePath=$(echo "$package_path" | sed "s|$dir||g")
        mkdir -p advocateResult/$packageName-$fileName-$test_func
        $pathToFullWorkflowExecutor -a $pathToAdvocate -p $adjustedPackagePath -f $dir -tf $file -t $test_func &> advocateResult/$packageName-$fileName-$test_func/output.txt
        cp -r $package_path/advocateTrace advocateResult/$packageName-$fileName-$test_func
        cp $package_path/results_machine.log advocateResult/$packageName-$fileName-$test_func
        cp $package_path/results_readable.log advocateResult/$packageName-$fileName-$test_func
        cp $package_path/times.log advocateResult/$packageName-$fileName-$test_func
        cp -r $package_path/rewritten_trace* advocateResult/$packageName-$fileName-$test_func 2>/dev/null
    done
    current_file=$((current_file+1))
done
