#!/bin/bash
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -a | --advocate)
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
# intialize variables
pathToPatchedGoRuntime="$pathToAdvocate/go-patch/bin/go"
pathToGoRoot="$pathToAdvocate/go-patch"
pathToOverheadInserter="$pathToAdvocate/toolchain/unitTestOverheadInserter/unitTestOverheadInserter"
pathToOverheadRemover="$pathToAdvocate/toolchain/unitTestOverheadRemover/unitTestOverheadRemover"

cd "$dir"
rm -r "advocateResults"
mkdir "advocateResults"
echo  "In directory: $dir"

export GOROOT=$pathToGoRoot
echo "Goroot exported"
test_files=$(find "$dir" -name "*_test.go")
total_files=$(echo "$test_files" | wc -l)
current_file=1
#echo "Test files: $test_files"
for file in $test_files; do
    echo "Progress: $current_file/$total_files"
    echo "Processing file: $file"
    package_path=$(dirname "$file")
    echo "Package path: $package_path"
    test_functions=$(grep -oE "[a-zA-Z0-9_]+ *Test[a-zA-Z0-9_]*" $file | sed 's/ *\(t *\*testing\.T\)//' | sed 's/func //')
    echo "Test functions: $test_functions"
    for test_func in $test_functions; do
        echo "Running test: $test_func" in "$file of package $package_path"
        echo "$pathToOverheadInserter -f $file -t $test_func"
        $pathToOverheadInserter -f $file -t $test_func
        if [ $? -ne 0 ]; then
          # If the overhead inserter fails, skip the test
          echo "Overhead inserter failed for $file $test_func. Skipping test."
          continue
        fi
        echo "$pathToPatchedGoRuntime test -count=1 -run=$test_func $package_path"
        $pathToPatchedGoRuntime test -count=1 -run="$test_func" "$package_path"
        echo "$pathToOverheadRemover -f $file -t $test_func"
        $pathToOverheadRemover -f "$file" -t "$test_func"
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        echo "mkdir -p advocateResults/$packageName/$fileName/$test_func"
        mkdir -p "advocateResults/$packageName/$fileName/$test_func"
        echo "mv $package_path/advocateTrace advocateResults/$packageName/$fileName/$test_func/advocateTrace"
        mv "$package_path/advocateTrace" "advocateResults/$packageName/$fileName/$test_func/advocateTrace"
    done
    current_file=$((current_file+1))
done