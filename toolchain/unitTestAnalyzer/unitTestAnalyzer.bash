#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi


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
      pathToOverheadInserter="$2"
      shift
      shift
      ;;
    -r|--overhead-remover)
      pathToOverheadRemover="$2"
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

if [ -z "$pathToPatchedGoRuntime" ] || [ -z "$pathToGoRoot" ] || [ -z "$pathToOverheadInserter" ] || [ -z "$pathToOverheadRemover" ] || [ -z "$dir" ];then
  echo "Usage: $0 -p <patched-go-runtime> -g <go-root> -i <overhead-inserter> -r <overhead-remover> -f <folder>"
  exit 1
fi




cd "$dir"
rm -r "advocateResults"
mkdir "advocateResults"
echo  "In directory: $dir"

export GOROOT=$pathToGoRoot
echo "Goroot exported"
test_files=$(find "$dir" -name "*_test.go")
echo "Test files: $test_files"
for file in $test_files; do
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
        $pathToPatchedGoRuntime test -count=1 -run="$test_func" "$package_path"
        $pathToOverheadRemover -f "$file" -t "$test_func"
        packageName=$(basename "$package_path")
        fileName=$(basename "$file")
        mkdir -p "advocateResults/$packageName/$fileName/$test_func"
        mv "$package_path/advocateTrace" "advocateResults/$packageName/$fileName/$test_func/advocateTrace"
    done
done