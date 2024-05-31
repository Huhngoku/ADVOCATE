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
skipped_tests=0
attempted_tests=0
#echo "Test files: $test_files"
for file in $test_files; do
    echo "Progress: $current_file/$total_files"
    echo "Processing file: $file"
    package_path=$(dirname "$file")
    test_functions=$(grep -oE "[a-zA-Z0-9_]+ *Test[a-zA-Z0-9_]*" $file | sed 's/ *\(t *\*testing\.T\)//' | sed 's/func //')
    for test_func in $test_functions; do
        attempted_tests=$((attempted_tests+1))
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        echo "Running full workflow for test: $test_func in package: $package_path in file: $file"
        adjustedPackagePath=$(echo "$package_path" | sed "s|$dir||g")
        directoryName="advocateResult/file($current_file)-test($attempted_tests)-$fileName-$test_func"
        echo "Creating directory: $directoryName"
        mkdir -p $directoryName
        $pathToFullWorkflowExecutor -a $pathToAdvocate -p $adjustedPackagePath -f $dir -tf $file -t $test_func &> $directoryName/output.txt
        # check if the test failed
        if [ $? -ne 0 ]; then
            echo "Test failed, check output.txt for more information. Skipping..."
            skipped_tests=$((skipped_tests+1))
            continue
        fi
        cp -r $package_path/advocateTrace $directoryName
        cp $package_path/results_machine.log $directoryName
        cp $package_path/results_readable.log $directoryName
        cp $package_path/times.log $directoryName
        cp -r $package_path/rewritten_trace* $directoryName 2>/dev/null
    done
    current_file=$((current_file+1))
done
echo "Finished fullworkflow for all tests"
echo "Attempted tests: $attempted_tests"
echo "Skipped tests: $skipped_tests"
