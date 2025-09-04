import math

def fibonacci(n):
    if n <= 1:
        return n
    else:
        return fibonacci(n-1) + fibonacci(n-2)

if __name__ == "__main__":
    print("Enter a number to generate its corresponding Fibonacci number: ")
    num = int(input())
    result = fibonacci(num)
    print(f"The {num}th Fibonacci number is {result}")
