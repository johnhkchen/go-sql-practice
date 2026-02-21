# T-001-03: astro-project-init - Structure

## File Operations Summary

### New Directories
1. `frontend/` - Root directory for Astro project
2. `frontend/src/` - Source code directory
3. `frontend/src/pages/` - Astro pages directory

### New Files
1. `frontend/package.json` - NPM package configuration
2. `frontend/astro.config.mjs` - Astro framework configuration
3. `frontend/tsconfig.json` - TypeScript configuration
4. `frontend/.gitignore` - Frontend-specific Git ignore rules
5. `frontend/src/pages/index.astro` - Home page component

### Modified Files
1. `.gitignore` - Add frontend build artifacts and dependencies

## Detailed File Specifications

### `frontend/package.json`
**Purpose:** Define project metadata, dependencies, and scripts
**Type:** JSON configuration file
**Key fields:**
- `name`: "go-sql-practice-frontend"
- `type`: "module" (ESM support)
- `version`: "0.0.1"
- `private`: true
- `dependencies`: `{"astro": "^5.17.3"}`
- `engines`: `{"node": ">=24"}`
- `scripts`: dev, build, preview commands

### `frontend/astro.config.mjs`
**Purpose:** Configure Astro for static output
**Type:** JavaScript ES module
**Exports:** Astro configuration object
**Key settings:**
- `output`: 'static' (no SSR)
- `build.assets`: 'assets' (asset directory name)
- `build.format`: 'directory' (clean URLs)

### `frontend/tsconfig.json`
**Purpose:** TypeScript compiler configuration
**Type:** JSON configuration file
**Extends:** "astro/tsconfigs/base"
**Compiler options:**
- Inherits Astro's base configuration
- JSX support for potential React components

### `frontend/.gitignore`
**Purpose:** Exclude frontend-specific files from Git
**Type:** Plain text gitignore file
**Patterns:**
- `node_modules/` - NPM dependencies
- `dist/` - Build output
- `.env*` - Environment files
- Editor and OS files

### `frontend/src/pages/index.astro`
**Purpose:** Smoke test home page
**Type:** Astro component file
**Structure:**
- HTML5 doctype and structure
- `<h1>Link Bookmarks</h1>` heading
- Minimal CSS reset (optional)
- No JavaScript initially

### `.gitignore` (root)
**Purpose:** Ensure frontend artifacts excluded
**Type:** Plain text gitignore file
**Additions:**
- `frontend/node_modules/`
- `frontend/dist/`

## Directory Structure After Changes

```
go-sql-practice/
├── .flox/
│   └── env/
│       └── manifest.toml
├── .git/
├── .gitignore (modified)
├── CLAUDE.md
├── docs/
│   ├── active/
│   │   ├── stories/
│   │   ├── tickets/
│   │   └── work/
│   │       ├── T-001-00/
│   │       └── T-001-03/
│   │           ├── research.md
│   │           ├── design.md
│   │           └── structure.md
│   └── knowledge/
└── frontend/ (new)
    ├── .gitignore (new)
    ├── astro.config.mjs (new)
    ├── package.json (new)
    ├── tsconfig.json (new)
    └── src/ (new)
        └── pages/ (new)
            └── index.astro (new)
```

## Public Interfaces

### NPM Scripts Interface
**Commands exposed via package.json:**
- `npm run dev` - Start development server on port 4321
- `npm run build` - Build static site to `dist/`
- `npm run preview` - Preview production build

### Build Output Interface
**Location:** `frontend/dist/`
**Structure:**
```
dist/
├── index.html
├── assets/
│   └── [hashed CSS/JS files]
└── [other pages].html
```

### Development Server Interface
**Port:** 4321 (Astro default)
**URL:** http://localhost:4321
**Hot reload:** Enabled by default

## Internal Organization

### Astro Project Conventions
- Pages in `src/pages/` auto-generate routes
- Components would go in `src/components/` (future)
- Layouts would go in `src/layouts/` (future)
- Static assets in `public/` directory (future)

### Build Pipeline Integration
- Frontend builds independently via `npm run build`
- Output in `dist/` ready for Go embedding
- No server runtime dependencies

## Module Boundaries

### Frontend Module
**Responsibility:** UI rendering and static asset generation
**Dependencies:** Node.js, Astro
**Outputs:** Static HTML/CSS/JS files
**Inputs:** None (static for now)

### Go Backend Module (future)
**Responsibility:** Serve embedded frontend files
**Dependencies:** Go embed package
**Consumes:** `frontend/dist/` contents

## Configuration Boundaries

### Node/NPM Configuration
- `package.json` - Dependencies and scripts
- `package-lock.json` - Dependency lock (generated)
- `.nvmrc` (optional future) - Node version

### Astro Configuration
- `astro.config.mjs` - Framework settings
- `tsconfig.json` - TypeScript settings

### Git Configuration
- `.gitignore` files at appropriate levels
- No committed build artifacts or dependencies

## Change Ordering

1. **Create directory structure** - `frontend/src/pages/`
2. **Create configuration files** - package.json, astro.config.mjs, tsconfig.json
3. **Create .gitignore files** - Frontend-specific and root updates
4. **Create index.astro** - Smoke test page
5. **Install dependencies** - `npm install` to generate lock file
6. **Verify build** - `npm run build` to test configuration

## Validation Criteria

### Structural Validation
- [ ] `frontend/` directory exists
- [ ] All configuration files present
- [ ] Source structure follows Astro conventions

### Build Validation
- [ ] `npm install` completes without errors
- [ ] `npm run build` generates `dist/` directory
- [ ] `dist/index.html` contains "Link Bookmarks" text

### Integration Validation
- [ ] Git ignores `node_modules/` and `dist/`
- [ ] Build output suitable for Go embedding
- [ ] No server-side code in output

## Non-Changes

### Not Modified
- Go module files (don't exist yet)
- Makefile (doesn't exist yet)
- CI/CD configuration (doesn't exist yet)
- Backend code (doesn't exist yet)

### Not Created
- Additional Astro pages
- Component library
- Styling system
- Client-side routing
- API integration code

## Notes

The structure establishes a minimal but complete Astro project that satisfies all acceptance criteria. It provides a foundation for future enhancements while keeping the initial implementation focused and testable.