package handlers

import (
	"bolter/utils"
	"encoding/json"
	"io"
	"net/http"
)

// Response is a simple structure for JSON responses
type TemplateResponse struct {
	Prompts   []string `json:"promts"`
	UiPrompts []string `json:"uiPrompts"`
}

type TemplateErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type TemplateRequestBody struct {
	Prompt string `json:"prompt"`
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	json_response := `{
		"promts": [
			"For all designs I ask you to make, have them be beautiful, not cookie cutter. Make webpages that are fully featured and worthy for production.\n\nBy default, this template supports JSX syntax with Tailwind CSS classes, React hooks, and Lucide React for icons. Do not install other packages for UI themes, icons, etc unless absolutely necessary or I request them.\n\nUse icons from lucide-react for logos.\n\nUse stock photos from unsplash where appropriate, only valid URLs you know exist. Do not download the images, only link to them in image tags.\n",
			"Here is an artifact that contains all files of the project visible to you.\\nConsider the contents of ALL files in the project.\\n\\n\n\u003cboltArtifact id=\"project-import\" title=\"Project Files\"\u003e\u003cboltAction type=\"file\" filePath=\"eslint.config.js\"\u003e\nimport js from '@eslint/js';\nimport globals from 'globals';\nimport reactHooks from 'eslint-plugin-react-hooks';\nimport reactRefresh from 'eslint-plugin-react-refresh';\nimport tseslint from 'typescript-eslint';\n\nexport default tseslint.config(\n  { ignores: ['dist'] },\n  {\n    extends: [js.configs.recommended, ...tseslint.configs.recommended],\n    files: ['**/*.{ts,tsx}'],\n    languageOptions: {\n      ecmaVersion: 2020,\n      globals: globals.browser,\n    },\n    plugins: {\n      'react-hooks': reactHooks,\n      'react-refresh': reactRefresh,\n    },\n    rules: {\n      ...reactHooks.configs.recommended.rules,\n      'react-refresh/only-export-components': [\n        'warn',\n        { allowConstantExport: true },\n      ],\n    },\n  }\n);\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"index.html\"\u003e\u003c!doctype html\u003e\n\u003chtml lang=\"en\"\u003e\n  \u003chead\u003e\n    \u003cmeta charset=\"UTF-8\" /\u003e\n    \u003clink rel=\"icon\" type=\"image/svg+xml\" href=\"/vite.svg\" /\u003e\n    \u003cmeta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" /\u003e\n    \u003ctitle\u003eVite + React + TS\u003c/title\u003e\n  \u003c/head\u003e\n  \u003cbody\u003e\n    \u003cdiv id=\"root\"\u003e\u003c/div\u003e\n    \u003cscript type=\"module\" src=\"/src/main.tsx\"\u003e\u003c/script\u003e\n  \u003c/body\u003e\n\u003c/html\u003e\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"package.json\"\u003e{\n  \"name\": \"vite-react-typescript-starter\",\n  \"private\": true,\n  \"version\": \"0.0.0\",\n  \"type\": \"module\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\",\n    \"lint\": \"eslint .\",\n    \"preview\": \"vite preview\"\n  },\n  \"dependencies\": {\n    \"lucide-react\": \"^0.344.0\",\n    \"react\": \"^18.3.1\",\n    \"react-dom\": \"^18.3.1\"\n  },\n  \"devDependencies\": {\n    \"@eslint/js\": \"^9.9.1\",\n    \"@types/react\": \"^18.3.5\",\n    \"@types/react-dom\": \"^18.3.0\",\n    \"@vitejs/plugin-react\": \"^4.3.1\",\n    \"autoprefixer\": \"^10.4.18\",\n    \"eslint\": \"^9.9.1\",\n    \"eslint-plugin-react-hooks\": \"^5.1.0-rc.0\",\n    \"eslint-plugin-react-refresh\": \"^0.4.11\",\n    \"globals\": \"^15.9.0\",\n    \"postcss\": \"^8.4.35\",\n    \"tailwindcss\": \"^3.4.1\",\n    \"typescript\": \"^5.5.3\",\n    \"typescript-eslint\": \"^8.3.0\",\n    \"vite\": \"^5.4.2\"\n  }\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"postcss.config.js\"\u003eexport default {\n  plugins: {\n    tailwindcss: {},\n    autoprefixer: {},\n  },\n};\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tailwind.config.js\"\u003e/** @type {import('tailwindcss').Config} */\nexport default {\n  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],\n  theme: {\n    extend: {},\n  },\n  plugins: [],\n};\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.app.json\"\u003e{\n  \"compilerOptions\": {\n    \"target\": \"ES2020\",\n    \"useDefineForClassFields\": true,\n    \"lib\": [\"ES2020\", \"DOM\", \"DOM.Iterable\"],\n    \"module\": \"ESNext\",\n    \"skipLibCheck\": true,\n\n    /* Bundler mode */\n    \"moduleResolution\": \"bundler\",\n    \"allowImportingTsExtensions\": true,\n    \"isolatedModules\": true,\n    \"moduleDetection\": \"force\",\n    \"noEmit\": true,\n    \"jsx\": \"react-jsx\",\n\n    /* Linting */\n    \"strict\": true,\n    \"noUnusedLocals\": true,\n    \"noUnusedParameters\": true,\n    \"noFallthroughCasesInSwitch\": true\n  },\n  \"include\": [\"src\"]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.json\"\u003e{\n  \"files\": [],\n  \"references\": [\n    { \"path\": \"./tsconfig.app.json\" },\n    { \"path\": \"./tsconfig.node.json\" }\n  ]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.node.json\"\u003e{\n  \"compilerOptions\": {\n    \"target\": \"ES2022\",\n    \"lib\": [\"ES2023\"],\n    \"module\": \"ESNext\",\n    \"skipLibCheck\": true,\n\n    /* Bundler mode */\n    \"moduleResolution\": \"bundler\",\n    \"allowImportingTsExtensions\": true,\n    \"isolatedModules\": true,\n    \"moduleDetection\": \"force\",\n    \"noEmit\": true,\n\n    /* Linting */\n    \"strict\": true,\n    \"noUnusedLocals\": true,\n    \"noUnusedParameters\": true,\n    \"noFallthroughCasesInSwitch\": true\n  },\n  \"include\": [\"vite.config.ts\"]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"vite.config.ts\"\u003eimport { defineConfig } from 'vite';\nimport react from '@vitejs/plugin-react';\n\n// https://vitejs.dev/config/\nexport default defineConfig({\n  plugins: [react()],\n  optimizeDeps: {\n    exclude: ['lucide-react'],\n  },\n});\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/App.tsx\"\u003eimport React from 'react';\n\nfunction App() {\n  return (\n    \u003cdiv className=\"min-h-screen bg-gray-100 flex items-center justify-center\"\u003e\n      \u003cp\u003eStart prompting (or editing) to see magic happen :)\u003c/p\u003e\n    \u003c/div\u003e\n  );\n}\n\nexport default App;\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/index.css\"\u003e@tailwind base;\n@tailwind components;\n@tailwind utilities;\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/main.tsx\"\u003eimport { StrictMode } from 'react';\nimport { createRoot } from 'react-dom/client';\nimport App from './App.tsx';\nimport './index.css';\n\ncreateRoot(document.getElementById('root')!).render(\n  \u003cStrictMode\u003e\n    \u003cApp /\u003e\n  \u003c/StrictMode\u003e\n);\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/vite-env.d.ts\"\u003e/// \u003creference types=\"vite/client\" /\u003e\n\u003c/boltAction\u003e\u003c/boltArtifact\u003e\\n\\nHere is a list of files that exist on the file system but are not being shown to you:\\n\\n  - .gitignore\\n  - package-lock.json\\n"
		],
		"uiPrompts": [
			"\n\u003cboltArtifact id=\"project-import\" title=\"Project Files\"\u003e\u003cboltAction type=\"file\" filePath=\"eslint.config.js\"\u003e\nimport js from '@eslint/js';\nimport globals from 'globals';\nimport reactHooks from 'eslint-plugin-react-hooks';\nimport reactRefresh from 'eslint-plugin-react-refresh';\nimport tseslint from 'typescript-eslint';\n\nexport default tseslint.config(\n  { ignores: ['dist'] },\n  {\n    extends: [js.configs.recommended, ...tseslint.configs.recommended],\n    files: ['**/*.{ts,tsx}'],\n    languageOptions: {\n      ecmaVersion: 2020,\n      globals: globals.browser,\n    },\n    plugins: {\n      'react-hooks': reactHooks,\n      'react-refresh': reactRefresh,\n    },\n    rules: {\n      ...reactHooks.configs.recommended.rules,\n      'react-refresh/only-export-components': [\n        'warn',\n        { allowConstantExport: true },\n      ],\n    },\n  }\n);\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"index.html\"\u003e\u003c!doctype html\u003e\n\u003chtml lang=\"en\"\u003e\n  \u003chead\u003e\n    \u003cmeta charset=\"UTF-8\" /\u003e\n    \u003clink rel=\"icon\" type=\"image/svg+xml\" href=\"/vite.svg\" /\u003e\n    \u003cmeta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" /\u003e\n    \u003ctitle\u003eVite + React + TS\u003c/title\u003e\n  \u003c/head\u003e\n  \u003cbody\u003e\n    \u003cdiv id=\"root\"\u003e\u003c/div\u003e\n    \u003cscript type=\"module\" src=\"/src/main.tsx\"\u003e\u003c/script\u003e\n  \u003c/body\u003e\n\u003c/html\u003e\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"package.json\"\u003e{\n  \"name\": \"vite-react-typescript-starter\",\n  \"private\": true,\n  \"version\": \"0.0.0\",\n  \"type\": \"module\",\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"vite build\",\n    \"lint\": \"eslint .\",\n    \"preview\": \"vite preview\"\n  },\n  \"dependencies\": {\n    \"lucide-react\": \"^0.344.0\",\n    \"react\": \"^18.3.1\",\n    \"react-dom\": \"^18.3.1\"\n  },\n  \"devDependencies\": {\n    \"@eslint/js\": \"^9.9.1\",\n    \"@types/react\": \"^18.3.5\",\n    \"@types/react-dom\": \"^18.3.0\",\n    \"@vitejs/plugin-react\": \"^4.3.1\",\n    \"autoprefixer\": \"^10.4.18\",\n    \"eslint\": \"^9.9.1\",\n    \"eslint-plugin-react-hooks\": \"^5.1.0-rc.0\",\n    \"eslint-plugin-react-refresh\": \"^0.4.11\",\n    \"globals\": \"^15.9.0\",\n    \"postcss\": \"^8.4.35\",\n    \"tailwindcss\": \"^3.4.1\",\n    \"typescript\": \"^5.5.3\",\n    \"typescript-eslint\": \"^8.3.0\",\n    \"vite\": \"^5.4.2\"\n  }\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"postcss.config.js\"\u003eexport default {\n  plugins: {\n    tailwindcss: {},\n    autoprefixer: {},\n  },\n};\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tailwind.config.js\"\u003e/** @type {import('tailwindcss').Config} */\nexport default {\n  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],\n  theme: {\n    extend: {},\n  },\n  plugins: [],\n};\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.app.json\"\u003e{\n  \"compilerOptions\": {\n    \"target\": \"ES2020\",\n    \"useDefineForClassFields\": true,\n    \"lib\": [\"ES2020\", \"DOM\", \"DOM.Iterable\"],\n    \"module\": \"ESNext\",\n    \"skipLibCheck\": true,\n\n    /* Bundler mode */\n    \"moduleResolution\": \"bundler\",\n    \"allowImportingTsExtensions\": true,\n    \"isolatedModules\": true,\n    \"moduleDetection\": \"force\",\n    \"noEmit\": true,\n    \"jsx\": \"react-jsx\",\n\n    /* Linting */\n    \"strict\": true,\n    \"noUnusedLocals\": true,\n    \"noUnusedParameters\": true,\n    \"noFallthroughCasesInSwitch\": true\n  },\n  \"include\": [\"src\"]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.json\"\u003e{\n  \"files\": [],\n  \"references\": [\n    { \"path\": \"./tsconfig.app.json\" },\n    { \"path\": \"./tsconfig.node.json\" }\n  ]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"tsconfig.node.json\"\u003e{\n  \"compilerOptions\": {\n    \"target\": \"ES2022\",\n    \"lib\": [\"ES2023\"],\n    \"module\": \"ESNext\",\n    \"skipLibCheck\": true,\n\n    /* Bundler mode */\n    \"moduleResolution\": \"bundler\",\n    \"allowImportingTsExtensions\": true,\n    \"isolatedModules\": true,\n    \"moduleDetection\": \"force\",\n    \"noEmit\": true,\n\n    /* Linting */\n    \"strict\": true,\n    \"noUnusedLocals\": true,\n    \"noUnusedParameters\": true,\n    \"noFallthroughCasesInSwitch\": true\n  },\n  \"include\": [\"vite.config.ts\"]\n}\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"vite.config.ts\"\u003eimport { defineConfig } from 'vite';\nimport react from '@vitejs/plugin-react';\n\n// https://vitejs.dev/config/\nexport default defineConfig({\n  plugins: [react()],\n  optimizeDeps: {\n    exclude: ['lucide-react'],\n  },\n});\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/App.tsx\"\u003eimport React from 'react';\n\nfunction App() {\n  return (\n    \u003cdiv className=\"min-h-screen bg-gray-100 flex items-center justify-center\"\u003e\n      \u003cp\u003eStart prompting (or editing) to see magic happen :)\u003c/p\u003e\n    \u003c/div\u003e\n  );\n}\n\nexport default App;\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/index.css\"\u003e@tailwind base;\n@tailwind components;\n@tailwind utilities;\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/main.tsx\"\u003eimport { StrictMode } from 'react';\nimport { createRoot } from 'react-dom/client';\nimport App from './App.tsx';\nimport './index.css';\n\ncreateRoot(document.getElementById('root')!).render(\n  \u003cStrictMode\u003e\n    \u003cApp /\u003e\n  \u003c/StrictMode\u003e\n);\n\u003c/boltAction\u003e\u003cboltAction type=\"file\" filePath=\"src/vite-env.d.ts\"\u003e/// \u003creference types=\"vite/client\" /\u003e\n\u003c/boltAction\u003e\u003c/boltArtifact\u003e"
		]
	}`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json_response))
	return
	// extract the body from the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// defer the closing of the body
	defer r.Body.Close()

	var data TemplateRequestBody
	// var data map[string]interface{}

	// convert the json to struct
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		// http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Get the singleton client
	client := utils.GetOpenRouterClient()

	// Create messages
	messages := []utils.Message{
		utils.SystemMessage("Return either node or react based on what do you think this project should be. Only return a single word either 'node' or 'react'. Do not return anything extra."),
		utils.UserMessage(data.Prompt),
	}

	// Call the API
	resp, err := client.ChatCompletion("meta-llama/llama-3.1-8b-instruct:free", messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(resp.Choices) > 0 {
		tech := resp.Choices[0].Message.Content
		if tech == "react" {
			json.NewEncoder(w).Encode(TemplateResponse{
				Prompts:   []string{utils.BasePrompt, utils.GetFSPrompt(tech)},
				UiPrompts: []string{utils.GetTechStackPrompt(tech)},
			})
		} else if tech == "node" {
			json.NewEncoder(w).Encode(TemplateResponse{
				Prompts:   []string{utils.GetFSPrompt(tech)},
				UiPrompts: []string{utils.GetTechStackPrompt(tech)},
			})
		}
	} else {
		json.NewEncoder(w).Encode(TemplateErrorResponse{
			Error: "No response choices returned",
		})
	}
}
