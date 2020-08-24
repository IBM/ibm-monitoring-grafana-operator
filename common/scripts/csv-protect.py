#!/usr/bin/python3
#
# Copyright 2020 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import os
import fnmatch
import yaml
import pathlib
import subprocess
import json

try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper


def packagePathes():
    projectDir = pathlib.Path(os.path.abspath(__file__)).parent.parent.parent
    packageFilePattern = "*.package.yaml"
    packageFile = ""
    catalogDir = ""
    for currentFile in pathlib.Path(pathlib.Path.joinpath(projectDir, "deploy", "olm-catalog")).rglob(packageFilePattern):
        packageFile = currentFile
        catalogDir = pathlib.Path(packageFile).parent
        break
    return str(projectDir), str(catalogDir), str(packageFile)


def modifiedFiles(projectDir):
    process = subprocess.Popen(['git', 'diff', "--name-only"],
                               stdout=subprocess.PIPE,
                               universal_newlines=True,
                               cwd=str(projectDir))
    unstagedFiles = process.stdout.readlines()

    process = subprocess.Popen(['git', 'diff', "--name-only", "--cached"],
                               stdout=subprocess.PIPE,
                               universal_newlines=True,
                               cwd=str(projectDir))

    stagedFiles = process.stdout.readlines()

    changedFiles = unstagedFiles + stagedFiles
    for i in range(len(changedFiles)):
        changedFiles[i] = str(pathlib.Path(
            changedFiles[i].rstrip('\n')).absolute())
    return changedFiles


def devCSV(packageFile, catalogDir):
    devCSV = ""
    devCSVPath = ""
    isNew = True
    stream = open(packageFile, 'r')
    packageContent = yaml.load(stream, Loader=Loader)
    for channel in packageContent["channels"]:
        if channel["name"] == "dev":
            devCSV = channel["currentCSV"]
    for channel in packageContent["channels"]:
        if channel["name"] != "dev" and channel["currentCSV"] == devCSV:
            isNew = False

    csvFilePattern = "*.clusterserviceversion.yaml"
    for csv in pathlib.Path(str(catalogDir)).rglob(csvFilePattern):
        if devCSV in str(csv):
            devCSVPath = str(csv)

    return devCSV, devCSVPath, isNew


def allCSVs(catalogDir):
    csvFilePattern = "*.clusterserviceversion.yaml"
    csvNames = []
    for csv in pathlib.Path(catalogDir).rglob(csvFilePattern):
        csvNames.append(str(csv))
    return csvNames


def validateExampleCR(csv, catalogDir):
    print(("valiate CR examples defined in csv: {0}".format(
        str(pathlib.Path(csv).relative_to(catalogDir)))))
    stream = open(csv, 'r')
    csvContent = yaml.load(stream, Loader=Loader)
    crs = csvContent["metadata"]["annotations"]["alm-examples"]
    json.loads(crs)
    print("CR examples are validated")
    return True


def main():
    print("start to check csv files")
    projectDir, catalogDir, packageFile = packagePathes()
    print(("project dir: {0}".format(projectDir)))
    changedFiles = modifiedFiles(projectDir)
    devCSVName, devCSVPath, devCSVIsNew = devCSV(packageFile, catalogDir)
    print(("find dev csv: {0}".format(devCSVName)))

    if devCSVName == "":
        print("ERROR: dev channel is not defined")
        exit(1)
    csvs = allCSVs(catalogDir)
    for csv in csvs:
        for changeFile in changedFiles:
            if csv == changeFile:
                if csv == devCSVName and (not devCSVIsNew):
                    print(("ERROR: modifing csv: {0}".format(csv)))
                    exit(1)
                if csv != devCSVPath:
                    print(("ERROR: modifing csv: {0}".format(csv)))
                    exit(1)
                if not validateExampleCR(devCSVPath, catalogDir):
                    print((
                        "ERROR: failed to validate csv: {0}".format(devCSVPath)))
                    exit(1)
    print("csv check passed")


if __name__ == "__main__":
    main()
