# T-001-03: astro-project-init - Design

## Objective
Create an Astro frontend project configured for static output that can be embedded into the Go binary, with a minimal index page as a smoke test.

## Design Options

### Option 1: Astro CLI Scaffolding
**Approach:** Use `npm create astro@latest` with automated responses
- **Pros:**
  - Official Astro project structure
  - Includes all standard configurations
  - TypeScript setup included by default
  - Best practices baked in
- **Cons:**
  - May include unnecessary boilerplate
  - Interactive prompts need automation
  - Might pull latest Astro instead of 5.17.3
- **Verdict:** Rejected - Version control and automation complexity

### Option 2: Manual Minimal Setup
**Approach:** Create files manually with minimal configuration
- **Pros:**
  - Complete control over structure
  - No unnecessary dependencies
  - Exactly what's needed, nothing more
  - Predictable and reproducible
- **Cons:**
  - Risk of missing important configs
  - More manual work
  - Need to know exact requirements
- **Verdict:** Selected - Best control and predictability

### Option 3: Template Repository Clone
**Approach:** Clone an Astro starter template and modify
- **Pros:**
  - Quick setup with examples
  - Working configuration guaranteed
- **Cons:**
  - External dependency on template
  - May have outdated dependencies
  - Cleanup of unwanted features needed
- **Verdict:** Rejected - External dependency risk

## Selected Approach: Manual Minimal Setup

### Rationale
1. **Version Control**: Explicit control over Astro version (5.17.3)
2. **Simplicity**: Only required files, no bloat
3. **Reproducibility**: No external dependencies or interactive processes
4. **Clarity**: Every file has a clear purpose
5. **Integration**: Optimized for embedding in Go binary

## Implementation Design

### Directory Structure
```
frontend/
├── package.json           # Dependencies and scripts
├── astro.config.mjs       # Astro configuration
├── tsconfig.json          # TypeScript configuration
├── .gitignore            # Local git ignore for frontend
└── src/
    └── pages/
        └── index.astro    # Smoke test page
```

### Package.json Configuration
```json
{
  "name": "go-sql-practice-frontend",
  "type": "module",
  "version": "0.0.1",
  "private": true,
  "scripts": {
    "dev": "astro dev",
    "build": "astro build",
    "preview": "astro preview"
  },
  "dependencies": {
    "astro": "^5.17.3"
  },
  "engines": {
    "node": ">=24"
  }
}
```

### Astro Configuration
```javascript
// astro.config.mjs
import { defineConfig } from 'astro/config';

export default defineConfig({
  output: 'static',
  build: {
    // Ensure clean builds for embedding
    assets: 'assets',
    format: 'directory'
  }
});
```

### TypeScript Configuration
```json
{
  "extends": "astro/tsconfigs/base",
  "compilerOptions": {
    "jsx": "react-jsx",
    "jsxImportSource": "react"
  }
}
```

### Index Page Design
```astro
---
// src/pages/index.astro
---
<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Link Bookmarks</title>
</head>
<body>
  <h1>Link Bookmarks</h1>
</body>
</html>
```

## Build Process Design

### Development Workflow
1. `cd frontend && npm install` - Install dependencies
2. `npm run dev` - Start development server on port 4321
3. Make changes, see hot reload
4. `npm run build` - Generate static files in `dist/`

### Production Build
1. `npm run build` creates `frontend/dist/`
2. T-001-04 will embed `frontend/dist/` into Go binary
3. Single binary serves static files at root path

## Dependency Management

### NPM Dependencies
- **astro@^5.17.3**: Core framework (caret for patch updates only)
- No additional dependencies for minimal setup
- TypeScript support included in Astro

### Node Version
- Engines field enforces Node.js >=24
- Matches flox environment (24.13.0)
- Prevents version mismatch issues

## Git Configuration

### Frontend .gitignore
```
# Dependencies
node_modules/

# Build output
dist/

# Environment
.env
.env.local

# Editor
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db
```

### Root .gitignore Updates
Need to ensure root `.gitignore` includes:
- `frontend/node_modules/`
- `frontend/dist/`

## Testing Strategy

### Smoke Test Validation
1. Page renders without errors
2. "Link Bookmarks" heading visible
3. Build completes successfully
4. Output in `dist/` directory
5. Static files ready for embedding

### Build Verification
- Check `dist/index.html` exists
- Verify assets directory created
- Confirm no server-side code generated
- Validate file sizes reasonable for embedding

## Integration Considerations

### With Go Embedding (T-001-04)
- Clean `dist/` structure for `embed.FS`
- Predictable file paths
- No server-side dependencies

### With Makefile (T-001-02)
- Simple `make frontend` target
- `cd frontend && npm run build`
- Clear success/failure indication

### With CI/CD (T-005-03)
- `npm ci` for reproducible installs
- Cache `node_modules/` for speed
- Build artifacts validation

## Risk Mitigation

### Version Pinning
- Explicit Astro version in package.json
- Node engine requirement
- Lockfile (`package-lock.json`) for reproducibility

### Build Failures
- Minimal configuration reduces failure points
- Clear error messages from Astro CLI
- Fallback to manual HTML if needed

### Embedding Issues
- Static output mode ensures embeddability
- No dynamic server requirements
- Clean directory structure

## Decision Summary
Manual minimal setup provides the best balance of simplicity, control, and integration with the Go backend. The approach ensures predictable builds, clean embedding, and maintains the exact version requirements specified in T-001-00.