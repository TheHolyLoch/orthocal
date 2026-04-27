## AI BEHAVIOR

This is our rule for general AI behavior including how an LLM should and should not react or handle operations.

1. **Verify Information**: Always verify information before presenting it. Do not make assumptions or speculate without clear evidence.

2. **File-by-File Changes**: Make changes file by file & give me a chance to spot mistakes.

3. **No Apologies**: Never use apologies.

4. **No Understanding Feedback**: Avoid giving feedback about understanding in comments or documentation.

5. **No Summaries**: Don't summarize changes made.

6. **No Inventions**: Don't invent changes other than what's explicitly requested.

7. **No Unnecessary Updates**: Don't suggest updates or changes to files when there are no actual modifications needed.

8. **Clarify Before Acting**: Do not make changes or proceed with any action unless you are 100% clear on the user's intentions and meaning. Ask for clarification if there is any ambiguity.

9. **Use Tabs for Indentation**: Always use tabs instead of spaces for indentation in code.

10. **No Sudo Access**: Never use `sudo` commands. If a command or script requires `sudo` or system package installation (like `apt install`), provide the command to the user and instruct them to run it manually.

---

## AI Response Optimization

This is our rule for how the AI should respond in chat conversations. Use concise, information-dense phrasing. Prefer shorter word choices when meaning is clear.

**Important:** This rule applies only to LLM chat responses to the user. It does NOT apply to:
- README files or project documentation
- Code comments
- Markdown files or other written content

Chat in this concise style normally unless the user requests a more verbose response.

Avoid unnecessary repetition or filler language.

### Examples

```
Verbose: "In the event that you need to…"
Optimized: "If you need to…"
```

```
Verbose: "In accordance with"
Optimized: "Per" or "Following"
```

```
Verbose: "Due to the fact that"
Optimized: "Because"
```

```
Verbose: "At this point in time"
Optimized: "Now"
```

```
Verbose: "Make a decision regarding"
Optimized: "Decide"
```

```
Verbose: "In order to achieve"
Optimized: "To achieve"
```

```
Verbose: "Has the ability to"
Optimized: "Can"
```

```
Verbose: "It is important to note that"
Optimized: "Note that" or just omit
```

```
Verbose: "For the purpose of"
Optimized: "For" or "To"
```

```
Verbose: "A number of"
Optimized: "Several" or "Some"
```


---

## DIRECTORY STRUCTURE

This is our rules for project directory structures. This serves as an example for guidance on how a projects files should be organized.

###### Example project directory tree:
```
company_project_scraper/
	.env
	.env.example
	.gitignore
	assets/
		company_metadata.json
		employees.sql
	core/
		config.py
		ftp.py
		irc.py
		events.py
		utils.py
	logs/
		projectname_yyyy-mm-dd.log
	Dockerfile
	LICENSE
	main.py
	requirements.txt
	README.md
	setup.sh
	test.py
```

###### Files descriptions:
- `.env` — secret environment variables, keep out of version control
- `.env.example` — template showing all required environment variables  
- `.gitignore` — files and directories to exclude from Git  
- `assets/` — static or downloaded data, e.g., CSVs, JSONs, SQL files; can be gitignored if large or generated locally  
- `core/` — reusable Python modules *(config, FTP, IRC, events, utility functions)*  
- `logs/` — application log files, typically dated  
- `Dockerfile` — container setup instructions  
- `LICENSE` — project license  
- `main.py` — main script that runs the project  
- `requirements.txt` — Python dependencies  
- `README.md` — project overview and instructions  
- `setup.sh` — setup script for OS packages, pip dependencies, Docker, etc.
- `test.py`— Unit testing or proof of example script to validate functionality where needed.

Not all of these files are required for every project, this is just an example as guidance.

All `.sh` bash files should be `chmod +x` .

---

## LICENSE

This is our rules for the LICENSE file that all projects will include.

###### LICENSE file contents:

```
ISC License

Copyright (c) {YYYY}, dgm <dgm@tuta.com>

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
```

**Note:** `{YYYY}` is replaced by the current year *(example: `2026`)*


---

## FILE HEADINGS

This is our rules for the file headings added to most files in the project directory. 

### Use proper shebangs:
| Language | Shebang                  |
| -------- | ------------------------ |
| Bash     | `#!/usr/bin/env bash`    |
| Python   | `#!/usr/bin/env python3` |

### Project name, credits, & file reference:

Any file that supports comments must include the following header comments at the top:

```text
# Project Name Here - Developed by dgm (dgm@tuta.com)
# project_dir_name/filename.ext
```

This includes Python, Bash, Dockerfile, `.env`, `.ini`, C, SQL, or any other text-based files, NOT markdown files.

This is placed immediately after the shebang, if one exists.

Use language specific comment prefixes if applicable like how `C` programming uses `//` for comments.


---

## GITIGNORE

This is our rules for creating and managing a `.gitignore` for a project to hide files you don't want in version control.

###### Examples:
- Sensitive settings or keys: `.env`
- Logs or temporary files: `logs/`
- Python cache: `__pycache__/`

Only include what is truly needed, this is not a file to over fluff.

Goal: keep secrets, logs, and auto-generated files out of Git.

---

## README

This is our rules for creating and managing a projects `README.md` file.

### Template
```
# Project Name

1–3 concise paragraphs explaining what the project does.

## Table of Contents
Clickable table of contents

## Requirements
- [Python](https://python.org)
	  - [request](https://pypi.org/project/requests)
- [Terraform](https://terraform.com)
  
###### System Recommendations
- 2GB Ram
- 50gb Disk

## Setup
pip install -r requirements or pip install [repo name]
or manual with git clone seup.sh etc.
or running setup.sh

talk about .env or settings to changeuse a table to give detail able them.

## Usage
how to run it the cli arguments in a table, etc

---

###### Mirrors: 
```

- Clear, direct, human tone, not hyper-intelligent language that is obviously AI generated.
- Never use emojis or the `—` character in any comments, code, or README files.
- Vertically align markdown tables based on the largest item in each column with only 1 space after it.
- Text in parentheses should be italicized using `*()* ` format in markdown.

- If a project is controversial, for example, to spam/flood something, a malicious proof of concept, or some kind of recon/intel data scraping, then include a disclaimer saying you are not liable, made for testing with legal permission or on your own servers.


---

## ENV

This is our rules for handling environment variables usage including requiring them, using them for sensitive settings, authentication, etc.

Create a `.env.example` file for any project that uses `.env` files. The example file will not be in `.gitignore`' and will serve as a template for whats expected in the`.env` file. Most people will do `cp .env.example .env` and then edit the `.env` to get started. 

### Bash

Basic one-liner example to prefix a bash script to require and load a local `.env` file:
```bash
[ -f .env ] && source .env || { echo "error: missing .env file" && exit 1 }
```

Here is an example of doing an undefined environment variable fallback:
```bash
echo ${ADMIN_USERNAME:-admin}
```

### Python:

Here is an example of how most `config.py` files will look. This example makes the `.env` file itself optional as environment variables may be set externally, via `systemd`, `docker`, etc:
```python
import os

# Load .env if it exists
if os.path.exists('.env'):
	try:
		from dotenv import load_dotenv
	except ImportError:
		raise ImportError('missing python-dotenv library (pip install python-dotenv)')
	else:
		load_dotenv()

# Authentication settings		
ADMIN_USERNAME = os.getenv('ADMIN_USERNAME', 'admin') # Has a default if undefined
ADMIN_PASSWORD = os.environ['API_KEY'] # Required

# Other settings
ENVIRONMENT    = 'prod'
VERSION        = '1.0.2b'
```

---

## PYTHON - General Guidelines

- Always use Python virtual environments (venvs) for all projects to ensure dependency isolation.
- Prefer standard Python libraries before third-party PyPI libraries.
- Comments in code should be informative for simply explaining the flow of the code in English.
- Inline comments are only for small clarifications, avoid over-commenting
- Prefer an asynchronous design with `asyncio`.
- Avoid using the `requests` library if possible. More often than not we will opt for `aiohttp` but there is almost always a standard library way to do anything `requstes` can do, meaning we don't have to depend on PyPI packages if we don't need to.
- Use `_` inside of large int's, like `250_000` instead of `250000` for readability.

- Avoid doing any of the following:
	- pointless variables declared
	- gigantic functions with too many lines of code
	- functions with too many arguments
	- adding un-necessary or un-asked for delays with `sleep`
	- using `subprocess` or `os` for something redundant, that can be done in Python *(like `curl` or `wget`)*.

- When killing a Python process that was started by the LLM, always verify afterwards that the process was actually terminated.
- For Flask applications, use separate files for web assets: `index.html`, `style.css`, and `utils.js` instead of inline code.


---

## PYTHON - Import Organization

Imports are done categorically, in order, with a blank line between each category:

1. **Standard Library - direct imports**
```python
import random
import sys
import time
```

2. **Standard Library - `from x import y`**
```python
from datetime import datetime
```

3. **Third-Party PyPI Libraries**
- All PyPI imports must match this format:
```python
try:
	import requests
except ImportError:
	raise ImportError('missing requests library (pip install requests)')
```

4. **Local Project Imports**
```python
from core import config, ftp, irc, utils
```

- Imports always belong at the top, never in the middle of code. The only exception is ones found after `if __name__ == '__main__:'.


---

## PYTHON - Function Conventions

- Use a naming convention format like `def event_message():` rather than `def eventMessage():` or other styling types.
- All functions should be in ABC order.
- Use type hints on arguments and return *(`def example(name: str) -> bool:`)*
- **Never** use `-> None:` as a return hint.
- Do not use `from typing import ...`, it makes type hints and return hints looks hideous.
- All functions should have a docstring matching the formats below for functions with and without arguments:

### Functions with Arguments
```python
def example_function(name: str, age: int) -> bool:
	'''
	This is a simple to the point minimal function with arguments description in 1 line.
	
	:param name: The customers name
	:param age: The customers age
	'''
```

### Functions without Arguments
```python
def simple_function():
	'''This is a simple to the point minimal function description in 1 line.```
```


---

## PYTHON - Code Style & Alignment

- Use single quotes only, double quotes inside nested f-strings.
- Vertically align assignments, dicts, parameters, comments, etc for readability, example:
```python
self.chunk_max  = args.chunk_max * 1024 * 1024 # MB
self.chunk_size = args.chunk_size
self.es_index   = args.index
```
- Inline comments should only have a single space before the `#`, not 2 spaces.
- Three blank lines before `if __name__ == '__main__'.
- Two blank lines after the last import.
- 1 blank line after a functions docstring line.
- 2 blank lines in between every function.


---

## PYTHON - Memory-Safe Data Handling

When reading large files or fetching large amounts of data *(files, URLs, streams)*, process incrementally to avoid high memory usage:
- `yield` or async generators
- Read line by line or in chunks

Storing a massive amount of data in variables can yield high memory usage which may be avoidable in certain situations.

Apply only when applicable *(don't overcomplicate small data handling)*.


---

## PYTHON - Logging with APV

The [APV](https://pypi.org/project/apv/) library is preferred in all Python projects for logging. It wraps the standard logging library, allowing you to use `logging.info`, `logging.error`, `logging.debug`, `logging.warn`, `logging.fatal`, etc throughout the project without defining a `logger` variable.

### Setup Example

```python
import argparse
import apv

# Setup parser
parser = argparse.ArgumentParser()
parser.add_argument('-d', '--debug', action='store_true', help='Enable debug logging')
args = parser.parse_args()

# Setup logging via apv
if args.debug:
	apv.setup_logging(level='DEBUG', log_to_disk=True, max_log_size=10*1024*1024, max_backups=5, log_file_name='application_log')
else:
	apv.setup_logging(level='INFO')
```

### Usage Throughout Code

Once APV is initialized, use standard logging calls anywhere in your project:

```python
import logging

logging.info('Application started')
logging.debug('Debug information here')
logging.error('An error occurred')
logging.warn('Warning message')
logging.fatal('Critical error')
```


---

## PYTHON - Project Code Examples

### Example main.py

```python
import asyncio
import logging


async def main():
	'''Main asynchronous entrypoint of the application'''

	# Start the main loop
	while True:
		try:
			pass # Your async code goes here

		except KeyboardInterrupt:
			logging.debug('Keyboard interrupt detected, closing connections & exiting...')

			# Exit the matrix
			break

		except Exception as e:
			logging.fatal(f'Critical connection error: {str(e)}')

		finally:
			await asyncio.sleep(15) # Delay before trying restarting the loop



if __name__ == '__main__':
	import argparse
	import apv

	# Setup parser
	parser = argparse.ArgumentParser()
	parser.add_argument('-d', '--debug', action='store_true', help='Enable debug logging')
	args = parser.parse_args()

	# Setup logging via apv
	if args.debug:
		apv.setup_logging(level='DEBUG', log_to_disk=True, max_log_size=100*1024*1024, max_backups=5, log_file_name='application_log')
	else:
		apv.setup_logging(level='INFO')

	logging.info('Application Name has started')

	# Run the async main
	asyncio.run(main())
```

---

### Example setup.sh
```bash
# Load environment variables
[ -f .env ] && source .env || { echo "Error: .env file not found"; exit 1; }

# Set xtrace, exit on error, & verbose mode (after loading environment variables)
set -xev

# Remove existing docker container if it exists
docker rm -f container_name 2>/dev/null || true

# Build the Docker image
docker build -t container_name .

# Run the Docker container with environment variables
docker run -d --name container_name \
  --restart unless-stopped \
  --network host \
  --hostname $(hostname) \
  -v /etc/os-release:/etc/os-release:ro \
  -e ES_HOST="${ES_HOST}" \
  -e ES_USERNAME="${ES_USERNAME}" \
  -e ES_PASSWORD="${ES_PASSWORD}" \
  -e ES_API_KEY="${ES_API_KEY}" \
  container_name
```

### Example Dockerfile
```Dockerfile
# Use the slim version of the Python image
FROM python:3-slim

# Create the app directory
RUN mkdir -p /app

# Set up in the application directory
WORKDIR /app

# Copy python requirements file
COPY requirements.txt .

# Set up Python environment and install dependencies
RUN python3 -m pip install --upgrade pip && python3 -m pip install --no-cache-dir --only-binary :all: -r requirements.txt --upgrade

# Cleanup the python requirements file (not needed at runtime)
RUN rm requirements.txt

# Copy only the necessary application files
COPY . .

# Create entrypoint script with restart capability
RUN printf '#!/bin/sh\ncd /app\nwhile true; do sleep 10 && python3 main.py -d; done' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
```


---

