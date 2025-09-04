## [Planner]


Here is a breakdown of the steps to build a CLI Fibonacci generator in Python:

Step 1: Install Python and create a new project directory

* Open a terminal or command prompt and install Python if not already installed.
* Create a new project directory for your Fibonacci generator using `mkdir` command.
* Change into the new directory using `cd` command.

Step 2: Set up the virtual environment (optional)

* If you are working on a team or want to maintain a consistent development environment, consider setting up a virtual environment for your project. This will create an isolated Python installation that you can use to develop and test your code without affecting the system Python installation.
* To set up a virtual environment, run `python -m venv .` in your terminal or command prompt. This will create a new virtual environment in the current directory.
* Activate the virtual environment by running `source bin/activate` (for Linux and macOS) or `venv\Scripts\activate` (for Windows).

Step 3: Install required packages

* Fibonacci numbers are a classic example of a recursive sequence, which means that each number is generated based on the previous two numbers. To generate Fibonacci numbers in Python, we need to import a recursive function from a library.
* You can use the `math` module or third-party libraries like `fibonacci` to generate Fibonacci numbers.
* Install the required packages by running `pip install math` (for Python 3) or `pip3 install math` (for Python 2).

Step 4: Write the code

* Write a function that takes an integer input and returns its corresponding Fibonacci number. The function should recursively call itself with the previous two numbers until it reaches the base case, which is the first two terms in the sequence.
* Here's an example implementation of the `fibonacci` function:
```
def fibonacci(n):
    if n <= 1:
        return n
    else:
        return fibonacci(n-1) + fibonacci(n-2)
```
Step 5: Test the code

* Test your code by running it in the terminal or command prompt. You can also use a testing framework like `unittest` to write automated tests for your code.
* Run the program with sample inputs and check that the output is correct.

Step 6: Package the code (optional)

* If you want to distribute your Fibonacci generator as a standalone program, you can package it in a format that can be easily installed and executed on other systems.
* You can use tools like `pyinstaller` or `py2exe` to create an executable for your program.
* Package the code by running `python -m pyinstaller --onefile fibonacci_generator.py` (for Python 3) or `python -m pyinstaller --onefile fibonacci_generator.py` (for Python 2).

Step 7: Deploy the program (optional)

* Once you have packaged your program, you can deploy it to a hosting platform like GitHub Pages or Heroku.
* You can also package and distribute your program through PyPI by running `python setup.py sdist` and then uploading the resulting distribution file to PyPI using the `twine` command.

And that's it! You have now built a CLI Fibonacci generator in Python that can be used to generate Fibonacci numbers based on user input.

## [Architect]


Design:
The design of this CLI Fibonacci generator is simple and straightforward. The program takes an integer input from the user and returns its corresponding Fibonacci number using a recursive function. The function recursively calls itself with the previous two numbers until it reaches the base case, which is the first two terms in the sequence.

Structure:
The structure of this CLI Fibonacci generator consists of three main components:

1. The `fibonacci` function, which takes an integer input and returns its corresponding Fibonacci number using a recursive approach.
2. The `main` function, which is the entry point of the program and calls the `fibonacci` function with user input as an argument.
3. The `if __name__ == "__main__":` block, which runs the `main` function when the program is executed.

Edge Cases:
Here are some edge cases that may arise while using this CLI Fibonacci generator:

1. Invalid input: If the user enters a non-integer value or a negative number, the program will raise a `ValueError` and exit with an error message.
2. Large inputs: The program may take longer to generate large Fibonacci numbers due to the recursive approach.
3. Incorrect output: If the user enters a number that is not in the sequence of Fibonacci numbers, the program will return incorrect results.

In conclusion, this CLI Fibonacci generator is a simple and effective way to generate Fibonacci numbers using Python. The program's design, structure, and edge cases are well-defined, making it easy to use and understand.

