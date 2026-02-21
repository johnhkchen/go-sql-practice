# T-001-03: astro-project-init - Plan

## Implementation Steps

### Step 1: Create Frontend Directory Structure
**Action:** Create the `frontend/` directory and subdirectories
**Commands:**
```bash
mkdir -p frontend/src/pages
```
**Verification:**
- Directory `frontend/src/pages/` exists
- Clean directory structure ready for files

### Step 2: Create package.json
**Action:** Create NPM package configuration with Astro dependency
**File:** `frontend/package.json`
**Content:**
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
**Verification:**
- File exists and is valid JSON
- Astro version is ^5.17.3
- Node engine requirement present

### Step 3: Create Astro Configuration
**Action:** Configure Astro for static output
**File:** `frontend/astro.config.mjs`
**Content:**
```javascript
import { defineConfig } from 'astro/config';

export default defineConfig({
  output: 'static',
  build: {
    assets: 'assets',
    format: 'directory'
  }
});
```
**Verification:**
- File exists and is valid JavaScript
- Static output mode configured

### Step 4: Create TypeScript Configuration
**Action:** Set up TypeScript for Astro components
**File:** `frontend/tsconfig.json`
**Content:**
```json
{
  "extends": "astro/tsconfigs/base",
  "compilerOptions": {
    "jsx": "react-jsx",
    "jsxImportSource": "react"
  }
}
```
**Verification:**
- File exists and is valid JSON
- Extends Astro base config

### Step 5: Create Frontend .gitignore
**Action:** Add frontend-specific Git exclusions
**File:** `frontend/.gitignore`
**Content:**
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
**Verification:**
- File exists with proper patterns
- Includes node_modules and dist

### Step 6: Update Root .gitignore
**Action:** Ensure root gitignore excludes frontend artifacts
**File:** `.gitignore`
**Additions:**
```
# Frontend
frontend/node_modules/
frontend/dist/
```
**Verification:**
- Root .gitignore updated
- Frontend paths included

### Step 7: Create Index Page
**Action:** Create minimal Astro page with required heading
**File:** `frontend/src/pages/index.astro`
**Content:**
```astro
---
// Minimal Astro page for smoke test
---
<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="Link Bookmarks Application">
  <title>Link Bookmarks</title>
  <style>
    body {
      font-family: system-ui, -apple-system, sans-serif;
      line-height: 1.6;
      margin: 2rem;
      color: #333;
    }
    h1 {
      color: #111;
    }
  </style>
</head>
<body>
  <h1>Link Bookmarks</h1>
  <p>Welcome to the Link Bookmarks application.</p>
</body>
</html>
```
**Verification:**
- File exists at correct path
- Contains "Link Bookmarks" heading
- Valid Astro component syntax

### Step 8: Install Dependencies
**Action:** Run npm install in flox environment
**Commands:**
```bash
flox activate -- bash -c "cd frontend && npm install"
```
**Verification:**
- No npm errors
- `node_modules/` directory created
- `package-lock.json` generated
- Astro package installed at correct version

### Step 9: Test Development Server
**Action:** Verify dev server starts and serves the page
**Commands:**
```bash
flox activate -- bash -c "cd frontend && npm run dev" &
sleep 5
curl -s http://localhost:4321 | grep "Link Bookmarks"
kill %1
```
**Verification:**
- Dev server starts on port 4321
- Page loads with correct heading
- No console errors

### Step 10: Build Production Assets
**Action:** Generate static build output
**Commands:**
```bash
flox activate -- bash -c "cd frontend && npm run build"
```
**Verification:**
- Build completes without errors
- `frontend/dist/` directory created
- `frontend/dist/index.html` exists
- HTML contains "Link Bookmarks" heading

### Step 11: Verify Build Output
**Action:** Check that build output is suitable for embedding
**Commands:**
```bash
ls -la frontend/dist/
test -f frontend/dist/index.html && echo "Index exists"
grep -q "Link Bookmarks" frontend/dist/index.html && echo "Content verified"
```
**Verification:**
- Static files only (no server code)
- Clean directory structure
- Assets properly organized

### Step 12: Commit Changes
**Action:** Stage and commit all new files
**Commands:**
```bash
git add frontend/
git add .gitignore
git commit -m "feat: initialize Astro frontend with static configuration

- Create frontend/ directory with Astro 5.17.3
- Configure for static output mode
- Add minimal index page with 'Link Bookmarks' heading
- Set up build pipeline for embedding in Go binary
- Specify Node.js >=24 engine requirement"
```
**Verification:**
- All files committed
- No build artifacts in commit
- Clean git status

## Testing Strategy

### Unit Testing
Not applicable for initial setup - no business logic yet.

### Integration Testing
1. **Build Integration:** Verify `npm run build` produces embeddable output
2. **Dev Server:** Ensure development workflow functions
3. **Version Compatibility:** Confirm Node.js 24 and Astro 5.17.3 work together

### Smoke Testing
1. **Page Render:** Index page displays "Link Bookmarks" heading
2. **Build Success:** Production build completes without errors
3. **Static Output:** No server-side code in dist/

## Rollback Strategy

If any step fails:
1. Remove `frontend/` directory: `rm -rf frontend/`
2. Revert .gitignore changes: `git checkout -- .gitignore`
3. Clean git status: `git clean -fd`
4. Restart from Step 1 with fixes

## Dependencies

### External Dependencies
- Node.js 24.13.0 (via flox)
- NPM (bundled with Node.js)
- Astro 5.17.3 (via npm)

### Internal Dependencies
- T-001-00 completion (flox environment ready)
- Git repository initialized
- Write permissions in project root

## Success Criteria

### Required Outcomes
- [x] `frontend/` directory with valid Astro project
- [x] `astro.config.mjs` configured for static output
- [x] `package.json` with Node >=24 requirement
- [x] `npm run build` produces `frontend/dist/`
- [x] Index page renders "Link Bookmarks" heading
- [x] `.gitignore` excludes node_modules and dist

### Quality Checks
- No npm audit vulnerabilities (high/critical)
- Build time under 10 seconds
- Clean console output (no warnings)
- Git repository remains clean

## Notes

- Use flox activate for all npm commands to ensure correct Node version
- The index page is intentionally minimal - just a smoke test
- Future tickets will enhance the frontend with actual functionality
- Build output structure is critical for T-001-04 (embed-frontend)