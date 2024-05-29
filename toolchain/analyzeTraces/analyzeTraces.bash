#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: $0 <pathToAnalyzer> <directory>"
    exit 1
fi
if [ -z "$2" ]; then
    echo "Usage: $0 <pathToAnalyzer> <directory>"
    exit 1
fi
path_to_analyzer=$1
dir_path=$2
directories=$(find "$dir_path" -type d -name "advocateTrace")
num_dirs=$(echo "$directories" | wc -l)
current_dir=1
for dir in $directories; do
    echo "Processing directory $current_dir of $num_dirs"
    echo "Directory Name: $dir"
    $path_to_analyzer -t $dir
    #$path_to_analyzer -t $dir &>"$dir/../advocateAnalysis.txt"
    current_dir=$((current_dir+1))
done 