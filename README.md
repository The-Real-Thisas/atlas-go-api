# Atlas Revision API

This is the core API for the Atlas Revision project. It provides question data for the website. All questions are from CIE.

---
## Initialize Data

> Note: This is assuming that data.json is present and the commands are run from project root.

```bash
# Creates questions.json
cd data && python3 questions.py && cd ..

# Creates topics.py
cd data && python3 topics.py && cd ..

# One liner
cd data && python3 questions.py && python3 topics.py && cd ..
```

---
## Build Project

> Note: This assuming that the json files `questions.json` and `topics.json` is present and the command is being run from the project root.

```bash
./build.sh
```

---
## Test Project

> Note: This assuming that the json files `questions.json` and `topics.json` is present and the command is being run from the project root.

```bash
./test.sh
```