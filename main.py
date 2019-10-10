import argparse
from diff_parser import DiffParser
from ast_parser import ASTParser


def main():
    if part == "1":
        print("Running Part 1 - Parsing diff files from... diffs/" + filename + "\n")
        diff_parser = DiffParser()
        with open("diffs/" + filename) as file:
            result = diff_parser.parse(file)
        result.to_text()

    elif part == "2":
        print("Running Part 2 - AST ")
        ast_parser = ASTParser()
        with open("ast/" + filename) as file:
            result = ast_parser.parse(file)
        result.to_text()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('-p', '--part', required=True)
    parser.add_argument('-f', '--filename', required=True)
    args = parser.parse_args()
    part = args.part
    filename = args.filename
    main()
