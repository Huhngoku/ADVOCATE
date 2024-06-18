while [[ $# -gt 0 ]]; do
	key="$1"
	case $key in
	-a | --advocate)
		pathToAdvocate="$2"
		shift
		shift
		;;
	-f | --folder)
		dir="$2"
		shift
		shift
		;;
	-m | --modulemode)
		modulemode="$2"
		shift
		shift
		;;
	-t | --test-name)
		testName="$2"
		shift
		shift
		;;
	-p | --package)
		package="$2"
		shift
		shift
		;;
	-tf | --test-file)
		file="$2"
		shift
		shift
		;;
	*)
		shift
		;;
	esac
done

pathToPatchedGoRuntime="$pathToAdvocate/go-patch/bin/go"
pathToGoRoot="$pathToAdvocate/go-patch"
pathToOverheadInserter="$pathToAdvocate/toolchain/unitTestOverheadInserter/unitTestOverheadInserter"
pathToOverheadRemover="$pathToAdvocate/toolchain/unitTestOverheadRemover/unitTestOverheadRemover"
pathToAnalyzer="$pathToAdvocate/analyzer/analyzer"

if [ -z "$pathToAdvocate" ]; then
	echo "Path to advocate is empty"
	exit 1
fi
if [ -z "$dir" ]; then
	echo "Directory is empty"
	exit 1
fi
if [ -z "$testName" ]; then
	echo "Test name is empty"
	exit 1
fi
if [ -z "$package" ]; then
	echo "Package is empty"
	exit 1
fi
if [ -z "$file" ]; then
	echo "Test file is empty"
	exit 1
fi

cd "$dir"
echo "In directory: $dir"
export GOROOT=$pathToGoRoot
echo "Goroot exported"
echo "Remove Overhead just in case"
$pathToOverheadRemover -f $file -t $testName
echo "Add Overhead"
echo "$pathToOverheadInserter -f $file -t $testName"
$pathToOverheadInserter -f $file -t $testName
if [ $? -ne 0 ]; then
	echo "Error in adding overhead"
	exit 1
fi
echo "Run test"
if [ "$modulemode" == "true" ]; then
	echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package"
	$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package"
else
	echo "$pathToPatchedGoRuntime test -count=1 -run=$testName ./$package"
	$pathToPatchedGoRuntime test -count=1 -run=$testName "./$package"
fi
if [ $? -ne 0 ]; then
	echo "Remove Overhead"
	$pathToOverheadRemover -f $file -t $testName
	echo "Error in running test, therefor overhead removed and full workflow stopped."
	exit 1
fi
echo "Remove Overhead"
echo "$pathToOverheadRemover -f $file -t $testName"
$pathToOverheadRemover -f $file -t $testName
$pathToAnalyzer -t "$dir/$package/advocateTrace"
rewritten_traces=$(find "./$package" -type d -name "rewritten_trace*")
for trace in $rewritten_traces; do
	rtracenum=$(echo $trace | grep -o '[0-9]*$')
	echo "Apply reorder overhead"
	echo $pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum"
	$pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum"
	if [ "$modulemode" == "true" ]; then
		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package"
		$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	else
		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName ./$package"
		$pathToPatchedGoRuntime test -count=1 -run=$testName "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	fi
	echo "Remove reorder overhead"
	echo "$pathToOverheadRemover -f $file -t $testName"
	$pathToOverheadRemover -f $file -t $testName
done
unset GOROOT
