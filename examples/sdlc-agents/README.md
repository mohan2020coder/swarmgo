# Run the code using this
go run .


🚀 SDLC Crew is starting with model: codellama
2025/09/04 13:33:19 🗂️  FileMemoryProvider initialized at memory
2025/09/04 13:33:19 [Planner] loading memory...
2025/09/04 13:33:19 📂 [Planner] no memory file, starting fresh
2025/09/04 13:33:19 [Planner] processing request...
2025/09/04 13:36:17 [Planner] response received
2025/09/04 13:36:17 📂 [Planner] no memory file, starting fresh
2025/09/04 13:36:17 ✏️  [Planner] adding message: role=user content="Build a CLI Fibonacci generator iin Python."
2025/09/04 13:36:17 💾 [Planner] saved 1 messages to memory
2025/09/04 13:36:17 📖 [Planner] loaded 1 messages from memory
2025/09/04 13:36:17 ✏️  [Planner] adding message: role=assistant content="\nHere is a breakdown of the steps to build a CLI Fibonacci generator in Python:\n..."
2025/09/04 13:36:17 💾 [Planner] saved 2 messages to memory
2025/09/04 13:36:17 [Writer] DOCS.md created with section Planner
2025/09/04 13:36:17 [Architect] loading memory...
2025/09/04 13:36:17 📂 [Architect] no memory file, starting fresh
2025/09/04 13:36:17 [Architect] processing request...
2025/09/04 13:38:13 [Architect] response received
2025/09/04 13:38:13 📂 [Architect] no memory file, starting fresh
2025/09/04 13:38:13 ✏️  [Architect] adding message: role=user content="\nHere is a breakdown of the steps to build a CLI Fibonacci generator in Python:\n..."
2025/09/04 13:38:13 💾 [Architect] saved 1 messages to memory
2025/09/04 13:38:13 📖 [Architect] loaded 1 messages from memory
2025/09/04 13:38:13 ✏️  [Architect] adding message: role=assistant content="\nDesign:\nThe design of tthis CLI Fibonacci generator is simple and straightforwar..."
2025/09/04 13:38:13 💾 [Architect] saved 2 messages to memory
2025/09/04 13:38:13 [Writer] DOCS.md appended with section Architect
2025/09/04 13:38:13 [Coder] loading memory...
2025/09/04 13:38:13 📂 [Coder] no memory file, starting fresh
2025/09/04 13:38:13 [Coder] processing request...
2025/09/04 13:39:29 [Coder] response received
2025/09/04 13:39:29 📂 [Coder] no memory file, starting fresh
2025/09/04 13:39:29 ✏️  [Coder] adding message: role=user content="\nHere is a breakdown of the steps  to build a CLI Fibonacci generator in Python:\n..."
2025/09/04 13:39:29 💾 [Coder] saved 1 messages to memory
2025/09/04 13:39:29 📖 [Coder] loaded 1 messages from memory
2025/09/04 13:39:29 ✏️  [Coder] adding message: role=assistant content="```\nimport math\n\ndef fibonaacci(n):\n    if n <= 1:\n        return n\n    else:\n   ..."
2025/09/04 13:39:29 💾 [Coder] saved 2 messages to memory
2025/09/04 13:39:29 [Reviewer] loading memory...
2025/09/04 13:39:29 📂 [Reviewer] no memory file, starting fresh
2025/09/04 13:39:29 [Reviewer] processing request...
2025/09/04 13:40:44 [Reviewer] response received
2025/09/04 13:40:44 📂 [Reviewer] no memory file, starting fresh
2025/09/04 13:40:44 ✏️  [Reviewer] adding message: role=user content="```\nimport math\n\ndef fibonacci(n):\n    if n <= 1:\n        return n\n    else:\n   ..."
2025/09/04 13:40:44 💾 [Reviewer] saved 1 messages to memory
2025/09/04 13:40:44 📖 [Reviewer] loaded 1 messages from memory
2025/09/04 13:40:44 ✏️  [Reviewer] adding message: role=assistant content="\nThis code is quite close to being correct, but there are a few minor issues tha..."
2025/09/04 13:40:44 💾 [Reviewer] saved 2 messages to memory
2025/09/04 13:40:44 [Writer] CODE.md created with section Reviewer
2025/09/04 13:40:44 [Writer] Execution report written to TEST_REPORT.md
✅ SDLC process completed. Generated: DOCS.md, CODE.md, TEST_REPORT.md, and source file.