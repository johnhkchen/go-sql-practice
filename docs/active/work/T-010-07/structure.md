# Structure: T-010-07 - Fix Go Build Failures

## File Modifications

### Modified File: routes/presentations.go

**Current Import Block (lines 3-9):**
```go
import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)
```

**New Import Block (lines 3-8):**
```go
import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)
```

**Change**: Remove line 4 containing `"crypto/subtle"`
**Impact**: File shrinks by one line, all subsequent line numbers shift up by 1

### Modified File: routes/links_view.go

**Current Import Block (lines 3-10):**
```go
import (
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)
```

**New Import Block (lines 3-9):**
```go
import (
	"net/http"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)
```

**Change**: Remove line 5 containing `"strings"`
**Impact**: File shrinks by one line, all subsequent line numbers shift up by 1

### Modified File: routes/links_view_test.go

**Current Import Block (lines 3-10):**
```go
import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)
```

**New Import Block (lines 3-11):**
```go
import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)
```

**Change**: Add `"github.com/pocketbase/pocketbase"` import before the existing `core` import
**Impact**: File grows by one line, function at line 201 shifts to line 202

## Import Grouping Convention

Based on analysis of other files in the routes package:
- Standard library imports come first
- Blank line separator
- Third-party imports grouped together
- Project-specific imports may be in a third group

The modifications maintain this convention:
- Standard library imports remain grouped at top
- Third-party pocketbase imports remain grouped together
- Blank line separation is preserved

## Line Number Impact Analysis

### routes/presentations.go
- Lines 1-3: Unchanged
- Line 4: DELETED (`"crypto/subtle"`)
- Old lines 5-327 → New lines 4-326
- Total lines: 327 → 326

### routes/links_view.go
- Lines 1-4: Unchanged
- Line 5: DELETED (`"strings"`)
- Old lines 6-73 → New lines 5-72
- Total lines: 73 → 72

### routes/links_view_test.go
- Lines 1-9: Unchanged
- Line 10: NEW (`"github.com/pocketbase/pocketbase"`)
- Line 11: Existing `"github.com/pocketbase/pocketbase/core"`
- Old lines 10-218 → New lines 11-219
- Total lines: 218 → 219
- Function `createTestLinkWithID` moves from line 201 to line 202

## Function Reference Updates

No function reference updates needed because:
1. Removed imports were already unused
2. Added import provides type that's already referenced in the code
3. No function signatures change
4. No function bodies change

## Public API Surface

No changes to public API:
- All exported functions retain same signatures
- No new exports added
- No exports removed

## Test Impact

The test file modification enables compilation of existing test:
- `createTestLinkWithID` helper can now use `*pocketbase.PocketBase` type
- No test logic changes required
- Other test functions unaffected

## Compilation Order

Files can be modified in any order since:
- Each change is independent
- No inter-file dependencies for these specific imports
- Go compiler processes the entire package at once

Recommended order for clarity:
1. routes/presentations.go (remove unused import)
2. routes/links_view.go (remove unused import)
3. routes/links_view_test.go (add missing import)