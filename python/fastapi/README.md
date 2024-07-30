appstack-fastapi
====

# Overview
Compare different application technology stacks. 

In this case, using: 

1. Python
2. FastAPI
3. SQLite3

# Running

```bash
python -m venv .venv/$(basename "$PWD")
source .venv/$(basename "$PWD")/bin/activate
pip install -r requirements.txt
pip freeze > requirements.txt:w
```



# References

1. FastAPI + SQLite: https://www.geeksforgeeks.org/fastapi-sqlite-databases/
