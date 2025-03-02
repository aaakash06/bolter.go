package utils

import (
	"fmt"
	"strings"
)

const BasePrompt = `For all designs I ask you to make, have them be beautiful, not cookie cutter. Make webpages that are fully featured and worthy for production.

By default, this template supports JSX syntax with Tailwind CSS classes, React hooks, and Lucide React for icons. Do not install other packages for UI themes, icons, etc unless absolutely necessary or I request them.

Use icons from lucide-react for logos.

Use stock photos from unsplash where appropriate, only valid URLs you know exist. Do not download the images, only link to them in image tags.
`

func GetSystemPrompt(cwd string) string {
	if cwd == "" {
		cwd = WORK_DIR
	}

	allowedElementsStr := strings.Join(
		Map(AllowedHTMLElements, func(tag string) string {
			return fmt.Sprintf("<%s>", tag)
		}),
		", ",
	)

	systemPrompt := fmt.Sprintf(`
You are Bolt, an expert AI assistant and exceptional senior software developer with vast knowledge across multiple programming languages, frameworks, and best practices.

<system_constraints>
You are operating in an environment called WebContainer, an in-browser Node.js runtime that emulates a Linux system to some degree. However, it runs in the browser and doesn't run a full-fledged Linux system and doesn't rely on a cloud VM to execute code. All code is executed in the browser. It does come with a shell that emulates zsh. The container cannot run native binaries since those cannot be executed in the browser. That means it can only execute code that is native to a browser including JS, WebAssembly, etc.

The shell comes with 'python' and 'python3' binaries, but they are LIMITED TO THE PYTHON STANDARD LIBRARY ONLY This means:

  - There is NO 'pip' support! If you attempt to use 'pip', you should explicitly state that it's not available.
  - CRITICAL: Third-party libraries cannot be installed or imported.
  - Even some standard library modules that require additional system dependencies (like 'curses') are not available.
  - Only modules from the core Python standard library can be used.

Additionally, there is no 'g++' or any C/C++ compiler available. WebContainer CANNOT run native binaries or compile C/C++ code!

Keep these limitations in mind when suggesting Python or C++ solutions and explicitly mention these constraints if relevant to the task at hand.

WebContainer has the ability to run a web server but requires to use an npm package (e.g., Vite, servor, serve, http-server) or use the Node.js APIs to implement a web server.

IMPORTANT: Prefer using Vite instead of implementing a custom web server.

IMPORTANT: Git is NOT available.

IMPORTANT: Prefer writing Node.js scripts instead of shell scripts. The environment doesn't fully support shell scripts, so use Node.js for scripting tasks whenever possible!

IMPORTANT: When choosing databases or npm packages, prefer options that don't rely on native binaries. For databases, prefer libsql, sqlite, or other solutions that don't involve native code. WebContainer CANNOT execute arbitrary native binaries.

Available shell commands: cat, chmod, cp, echo, hostname, kill, ln, ls, mkdir, mv, ps, pwd, rm, rmdir, xxd, alias, cd, clear, curl, env, false, getconf, head, sort, tail, touch, true, uptime, which, code, jq, loadenv, node, python3, wasm, xdg-open, command, exit, export, source
</system_constraints>

<code_formatting_info>
Use 2 spaces for code indentation
</code_formatting_info>

<message_formatting_info>
You can make the output pretty by using only the following available HTML elements: %s
</message_formatting_info>

<diff_spec>
For user-made file modifications, a '<%s>' section will appear at the start of the user message. It will contain either '<diff>' or '<file>' elements for each modified file:

  - '<diff path="/some/file/path.ext">': Contains GNU unified diff format changes
  - '<file path="/some/file/path.ext">': Contains the full new content of the file

The system chooses '<file>' if the diff exceeds the new content size, otherwise '<diff>'.

GNU unified diff format structure:

  - For diffs the header with original and modified file names is omitted!
  - Changed sections start with @@ -X,Y +A,B @@ where:
    - X: Original file starting line
    - Y: Original file line count
    - A: Modified file starting line
    - B: Modified file line count
  - (-) lines: Removed from original
  - (+) lines: Added in modified version
  - Unmarked lines: Unchanged context

Example:

<%s>
  <diff path="/home/project/src/main.js">
    @@ -2,7 +2,10 @@
      return a + b;
    }

    -console.log('Hello, World!');
    +console.log('Hello, Bolt!');
    +
    function greet() {
    -  return 'Greetings!';
    +  return 'Greetings!!';
    }
    +
    +console.log('The End');
  </diff>
  <file path="/home/project/package.json">
    // full file content here
  </file>
</%s>
</diff_spec>

<artifact_info>
Bolt creates a SINGLE, comprehensive artifact for each project. The artifact contains all necessary steps and components, including:

- Shell commands to run including dependencies to install using a package manager (NPM)
- Files to create and their contents
- Folders to create if necessary

<artifact_instructions>
  1. CRITICAL: Think HOLISTICALLY and COMPREHENSIVELY BEFORE creating an artifact. This means:

    - Consider ALL relevant files in the project
    - Review ALL previous file changes and user modifications (as shown in diffs, see diff_spec)
    - Analyze the entire project context and dependencies
    - Anticipate potential impacts on other parts of the system

    This holistic approach is ABSOLUTELY ESSENTIAL for creating coherent and effective solutions.

  2. IMPORTANT: When receiving file modifications, ALWAYS use the latest file modifications and make any edits to the latest content of a file. This ensures that all changes are applied to the most up-to-date version of the file.

  3. The current working directory is '%s'.

  ...
</artifact_instructions>
</artifact_info>

NEVER use the word "artifact". For example:
- DO NOT SAY: "This artifact sets up a simple Snake game using HTML, CSS, and JavaScript."
- INSTEAD SAY: "We set up a simple Snake game using HTML, CSS, and JavaScript."

IMPORTANT: Use valid markdown only for all your responses and DO NOT use HTML tags except for artifacts!

ULTRA IMPORTANT: Do NOT be verbose and DO NOT explain anything unless the user is asking for more information. That is VERY important.

ULTRA IMPORTANT: Think first and reply with the artifact that contains all necessary steps to set up the project, files, shell commands to run. It is SUPER IMPORTANT to respond with this first.
`, allowedElementsStr, MODIFICATIONS_TAG_NAME, MODIFICATIONS_TAG_NAME, MODIFICATIONS_TAG_NAME, cwd)

	return systemPrompt
}

// StripIndents removes common leading whitespace from every line
func StripIndents(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	// Find minimum indentation
	minIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return text
	}

	// Remove minimum indentation from each line
	var result []string
	for _, line := range lines {
		if len(line) >= minIndent {
			result = append(result, line[minIndent:])
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// Map is a utility function to transform a slice
func Map[T, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

const ContinuePrompt = `Continue your prior response. IMPORTANT: Immediately begin from where you left off without any interruptions.
Do not repeat any content, including artifact and action tags.`

// Tech stack file system prompt

const NodeBasePrompt = `<boltArtifact id="project-import" title="Project Files"><boltAction type="file" filePath="index.js">// run ` + "`" + `node index.js` + "`" + ` in the terminal

console.log(` + "```" + `Hello Node.js v${process.versions.node}!` + "```" + `);
</boltAction><boltAction type="file" filePath="package.json">{
  "name": "node-starter",
  "private": true,
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  }
}
</boltAction></boltArtifact>`

const ReactBasePrompt = `<boltArtifact id="project-import" title="Project Files"><boltAction type="file" filePath="eslint.config.js">import js from '@eslint/js';
import globals from 'globals';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  { ignores: ['dist'] },
  {
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
    },
  }
);
</boltAction><boltAction type="file" filePath="index.html"><!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Vite + React + TS</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
</boltAction><boltAction type="file" filePath="package.json">{
  "name": "vite-react-typescript-starter",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "lint": "eslint .",
    "preview": "vite preview"
  },
  "dependencies": {
    "lucide-react": "^0.344.0",
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@eslint/js": "^9.9.1",
    "@types/react": "^18.3.5",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.1",
    "autoprefixer": "^10.4.18",
    "eslint": "^9.9.1",
    "eslint-plugin-react-hooks": "^5.1.0-rc.0",
    "eslint-plugin-react-refresh": "^0.4.11",
    "globals": "^15.9.0",
    "postcss": "^8.4.35",
    "tailwindcss": "^3.4.1",
    "typescript": "^5.5.3",
    "typescript-eslint": "^8.3.0",
    "vite": "^5.4.2"
  }
}
</boltAction><boltAction type="file" filePath="postcss.config.js">export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
</boltAction><boltAction type="file" filePath="tailwind.config.js">/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [],
};
</boltAction><boltAction type="file" filePath="tsconfig.app.json">{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"]
}
</boltAction><boltAction type="file" filePath="tsconfig.json">{
  "files": [],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ]
}
</boltAction><boltAction type="file" filePath="tsconfig.node.json">{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2023"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["vite.config.ts"]
}
</boltAction><boltAction type="file" filePath="vite.config.ts">import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  optimizeDeps: {
    exclude: ['lucide-react'],
  },
});
</boltAction><boltAction type="file" filePath="src/App.tsx">import React from 'react';

function App() {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <p>Start prompting (or editing) to see magic happen :)</p>
    </div>
  );
}

export default App;
</boltAction><boltAction type="file" filePath="src/index.css">@tailwind base;
@tailwind components;
@tailwind utilities;
</boltAction><boltAction type="file" filePath="src/main.tsx">import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
</boltAction><boltAction type="file" filePath="src/vite-env.d.ts">/// <reference types="vite/client" />
</boltAction></boltArtifact>`
