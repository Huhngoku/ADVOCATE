#!/usr/bin/python
import sys
from argparse import ArgumentParser
import os


MIN_PYTHON = (3, 10)
if sys.version_info < MIN_PYTHON:
    sys.exit("Python %s.%s or later is required.\n" % MIN_PYTHON)


def parseArgs() -> dict:
    """
    Parces the command line arguments given to the program.
    Returns:
        args: the arguments given to the program
    """
    parser = ArgumentParser()
    parser.add_argument("-o", "--output", help="name of the output file",
                        dest="output")
    parser.add_argument("-p", "--program", help="path to the analyzed program",
                        dest="program")
    parser.add_argument("-t", "--trace", help="path to the trace file",
                        dest="trace")
    parser.add_argument("-e", "--errors", help="path to the file containing " +
                        "the errors", dest="bugs")
    parser.add_argument("-r0", "--runtime0",
                        help="runtime without the modified runtime", dest="r0",
                        type=float)
    parser.add_argument("-r1", "--runtime1",
                        help="runtime with the modified runtime, but without" +
                        " trace creation", dest="r1", type=float)
    parser.add_argument("-r2", "--runtime2",
                        help="runtime with the modified runtime and with " +
                        "trace creation", dest="r2", type=float)
    parser.add_argument("-r3", "--runtime for analysis", dest="r3", type=float)
    return parser.parse_args()


def parseProgram(args: dict) -> dict:
    """
    Parces the program given by the path and measure the number of go files
    and the number of lines of code.
    Args:
        args: the arguments given to the program
    Returns:
        dict: the number of go files and the number of lines of code
    """
    if not args.program:
        return {}

    path = args.program

    totalFiles = 0
    totalLines = 0
    for root, _, files in os.walk(path):
        for file in files:
            # only look at go files
            if not file.endswith(".go"):
                continue
            totalFiles += 1
            with open(os.path.join(root, file), "r") as f:
                lines = f.readlines()
                totalLines += len(lines)

    return {"files": totalFiles, "lines": totalLines}


def parseTrace(args: dict) -> dict:
    if not args.trace:
        return {}

    path = args.trace
    res = {
        "numberRoutines": 0,

        "numberAtomics": 0,
        "numberAtomicOperations": 0,

        "numberChannels": 0,
        "numberChannelOperations": 0,  # includes executed select cases

        "numberMutexes": 0,
        "numberMutexOperations": 0,

        "numberOnce": 0,
        "numberOnceOperations": 0,

        "numberSelects": 0,
        # number of available cases (including default)
        "numberSelectCases": 0,
        "numberSelectChanOps": 0,  # number of exec channel operations in sel
        "numberSelectDefaults": 0,  # number of exec default cases in sel

        "numberWaitGroups": 0,
        "numberWaitGroupOperations": 0,
    }

    known = {
        "atomic": [],
        "channel": [],
        "mutex": [],
        "once": [],
        "waitgroup": []
    }

    with open(path, 'r') as f:
        for line in f.readlines():
            res["numberRoutines"] += 1
            for elem in line.split(";"):
                fields = elem.split(",")
                match fields[0]:
                    case "A":
                        res["numberAtomicOperations"] += 1
                        if fields[2] not in known["atomic"]:
                            res["numberAtomics"] += 1
                            known["atomic"].append(fields[2])
                    case "C":
                        res["numberChannelOperations"] += 1
                        if fields[3] not in known["channel"]:
                            res["numberChannels"] += 1
                            known["channel"].append(fields[3])
                    case "M":
                        res["numberMutexOperations"] += 1
                        if fields[3] not in known["mutex"]:
                            res["numberMutexes"] += 1
                            known["mutex"].append(fields[3])
                    case "O":
                        res["numberOnceOperations"] += 1
                        if fields[3] not in known["once"]:
                            res["numberOnce"] += 1
                            known["once"].append(fields[3])
                    case "S":
                        res["numberSelects"] += 1
                        cases = fields[4].split("~")
                        res["numberSelectCases"] += len(cases)
                        if cases[-1] == "D":
                            res["numberSelectDefaults"] += 1
                        else:
                            res["numberSelectChanOps"] += 1
                    case "W":
                        res["numberWaitGroupOperations"] += 1
                        if fields[3] not in known["waitgroup"]:
                            res["numberWaitGroups"] += 1
                            known["waitgroup"].append(fields[3])
                    case _:
                        pass
    return res


def parseRuntime(args: dict) -> dict:
    res = {}

    if args.r0:
        res["runtime0"] = args.r0
    if args.r1:
        res["runtime1"] = args.r1
    if args.r2:
        res["runtime2"] = args.r2
    if args.r3:
        res["runtime3"] = args.r3

    digits = 3

    if args.r0 and args.r1:
        res_value = max(0, args.r1 - args.r0)
        res["overhead1-value"] = round(res_value, digits)
        res_perc = max(0, (args.r1 - args.r0) / args.r0) * 100
        res["overhead1-perc"] = round(res_perc, digits)

    if args.r0 and args.r2:
        res_value = max(0, args.r2 - args.r0)
        res["overhead2-value"] = round(res_value, digits)
        res_perc = max(0, (args.r2 - args.r0) / args.r0) * 100
        res["overhead2-perc"] = round(res_perc, digits)

    return res


def writeProgramOverview(args: dict, info: dict) -> str:
    res = "## Program \n"

    if not args.program:
        res += "No program given.\n"
        return res

    res += "| Info | Value |\n| - | - |\n"
    res += "| Number of go files | " + str(info["files"]) + "|\n"
    res += "| Number of lines of code |" + str(info["lines"]) + "|\n"
    return res


def writeTraceOverview(args: dict, info: dict) -> str:
    res = "## Trace \n"

    if not args.trace:
        res += "No trace file given.\n"
        return res

    res += "| Info | Value |\n| - | - |\n"

    res += "| Number of routines | " + str(info["numberRoutines"]) + "|\n"

    res += "| Number of atomic variables | " + str(info["numberAtomics"]) + \
        "|\n"
    res += "| Number of atomic operations | " + \
        str(info["numberAtomicOperations"]) + "|\n"

    res += "| Number of channels | " + str(info["numberChannels"]) + "|\n"
    res += "| Number of channel operations | " + \
        str(info["numberChannelOperations"]) + "|\n"

    res += "| Number of mutexes | " + str(info["numberMutexes"]) + "|\n"
    res += "| Number of mutex operations | " + \
        str(info["numberMutexOperations"]) + "|\n"

    res += "| Number of once variables | " + str(info["numberOnce"]) + "|\n"
    res += "| Number of once operations | " + \
        str(info["numberOnceOperations"]) + "|\n"

    res += "| Number of selects | " + str(info["numberSelects"]) + "|\n"
    res += "| Number of select cases | " + \
        str(info["numberSelectCases"]) + "|\n"
    res += "| Number of executed select channel operations | " + \
        str(info["numberSelectChanOps"]) + "|\n"
    res += "| Number of executed select default cases | " + \
        str(info["numberSelectDefaults"]) + "|\n"

    res += "| Number of waitgroups | " + str(info["numberWaitGroups"]) + "|\n"
    res += "| Number of waitgroup operations | " + \
        str(info["numberWaitGroupOperations"]) + "|\n"

    return res


def writeRuntimeOverview(args: dict, info: dict) -> str:
    res = "## Runtime \n"

    if not args.r0 and not args.r1 and not args.r2 and not args.r3:
        res += "No runtime info given.\n"
        return res

    res += "| Info | Value |\n| - | - |\n"

    if args.r0:
        res += "| Runtime without modifications | " + \
            str(info["runtime0"]) + "|\n"
    if args.r1:
        res += "| Runtime with modified runtime | " + \
            str(info["runtime1"]) + "|\n"
    if args.r2:
        res += "| Runtime with modified runtime and trace creation | " + \
            str(info["runtime2"]) + "|\n"
    if args.r0 and args.r1:
        res += "| Overhead of modified runtime [s] | " + \
            str(info["overhead1-value"]) + "|\n"
        res += "| Overhead of modified runtime [\%] | " + \
            str(info["overhead1-perc"]) + "|\n"

    if args.r0 and args.r2:
        res += "| Overhead of modified runtime and trace creation [s] | " + \
            str(info["overhead2-value"]) + "|\n"
        res += "| Overhead of modified runtime and trace creation [\%] | " + \
            str(info["overhead2-perc"]) + "|\n"

    if args.r3:
        res += "| Runtime for analysis [s] | " + str(info["runtime3"]) + "|\n"

    return res


def writeResults(args: dict) -> str:
    if not args.bugs:
        return ""

    res = "## Found Results\n"
    with open(args.bugs, "r") as f:
        for line in f.readlines():
            res += line.replace("\n", "").replace(
                "\t", "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;") + "\\\n"

    return res


def createOverview(args: dict, info: dict) -> dict:
    ouput = "Overview.md"
    if args.output:
        if args.output.endswith(".md"):
            output = args.output
        else:
            output = args.output + ".md"
    elif args.program:
        output = args.program.split(os.sep)[-1] + ".md"

    name = output.split(".")[-2]

    with open(output, "w") as f:
        f.write("# " + name + "\n\n")
        f.write(writeProgramOverview(args, info))
        f.write(writeTraceOverview(args, info))
        f.write(writeRuntimeOverview(args, info))
        f.write(writeResults(args))


def main():
    args = parseArgs()
    res = {}
    res.update(parseProgram(args))
    res.update(parseTrace(args))
    res.update(parseRuntime(args))

    createOverview(args, res)


if __name__ == "__main__":
    main()
