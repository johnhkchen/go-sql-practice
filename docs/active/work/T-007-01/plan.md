# Plan: Presentations Collection (T-007-01)

## Implementation Steps

### Step 1: Add createPresentationsCollection Function
**File:** `migrations/collections.go`
**Location:** After line 196, before closing brace

**Actions:**
1. Add new function with signature `func createPresentationsCollection(txApp core.App) error`
2. Add existence check for presentations collection
3. Return nil if already exists

**Verification:**
- Function compiles without errors
- No duplicate function names

**Commit:** "feat: add presentations collection creation function stub"

### Step 2: Add Dependency Checks
**File:** `migrations/collections.go`
**Location:** Within createPresentationsCollection function

**Actions:**
1. Find sync_sessions collection, return error if missing
2. Find users collection, store reference if exists
3. Store collection IDs for later use

**Verification:**
- Error returned if sync_sessions missing
- No error if users missing (optional)

**Commit:** "feat: add dependency verification for presentations collection"

### Step 3: Create Collection and Add Basic Fields
**File:** `migrations/collections.go`
**Location:** Within createPresentationsCollection function

**Actions:**
1. Create new base collection named "presentations"
2. Add name field (TextField, required, 1-255 chars)
3. Add step_count field (NumberField, required, min 1, integer)

**Verification:**
- Collection object created successfully
- Required fields have correct constraints

**Commit:** "feat: add basic fields to presentations collection"

### Step 4: Add JSON and Relation Fields
**File:** `migrations/collections.go`
**Location:** Within createPresentationsCollection function

**Actions:**
1. Add step_labels field (JSONField, optional)
2. Add active_session relation to sync_sessions
3. Add created_by relation if users collection exists

**Verification:**
- JSON field accepts array structures
- Relations point to correct collection IDs

**Commit:** "feat: add JSON and relation fields to presentations collection"

### Step 5: Set API Rules
**File:** `migrations/collections.go`
**Location:** Within createPresentationsCollection function

**Actions:**
1. Set ListRule and ViewRule to "" (public)
2. Set CreateRule to "@request.auth.id != ''" (authenticated)
3. Set UpdateRule to "@request.auth.id = created_by" (owner)
4. Set DeleteRule to "@request.auth.id = created_by" (owner)

**Verification:**
- Rules compile as valid PocketBase expressions
- Public can read, authenticated can create, owners can modify

**Commit:** "feat: configure API access rules for presentations"

### Step 6: Save Collection and Handle Errors
**File:** `migrations/collections.go`
**Location:** Within createPresentationsCollection function

**Actions:**
1. Call txApp.Save(collection)
2. Wrap and return any errors
3. Return nil on success

**Verification:**
- Errors properly wrapped with context
- Success returns nil

**Commit:** "feat: complete presentations collection save logic"

### Step 7: Integrate with Main Migration
**File:** `migrations/collections.go`
**Location:** Line 146, within createCollections function

**Actions:**
1. Add call to createPresentationsCollection after sync_sessions
2. Ensure error propagation

**Verification:**
- Migration calls new function
- Errors halt migration process

**Commit:** "feat: integrate presentations collection into migration flow"

### Step 8: Test Migration
**Actions:**
1. Start the application
2. Verify presentations collection created
3. Check all fields exist with correct types
4. Verify API rules work as expected

**Verification:**
- No errors during startup
- Collection visible in PocketBase admin (if available)
- Can create/read via API

**Commit:** No commit needed (testing only)

## Testing Strategy

### Unit Testing
Not applicable - PocketBase migrations tested via integration

### Integration Testing

#### Test 1: Fresh Database
1. Start with empty database
2. Run application
3. Verify all collections created
4. Check presentations has all fields

#### Test 2: Existing Collections
1. Start with tags, links, sync_sessions existing
2. Run application
3. Verify presentations added without affecting others
4. Verify idempotent (run twice, no errors)

#### Test 3: Missing Dependencies
1. Remove sync_sessions collection
2. Run application
3. Verify error returned
4. Verify helpful error message

#### Test 4: API Access
1. Try to list presentations (should succeed)
2. Try to create without auth (should fail)
3. Create with auth (should succeed)
4. Update own presentation (should succeed)
5. Update another's presentation (should fail)

### Manual Verification

#### Field Types
- name accepts text, rejects empty
- step_count accepts integers ≥ 1, rejects 0 or decimals
- step_labels accepts JSON arrays like ["Step 1", "Step 2"]
- active_session accepts sync_session record ID
- created_by populated with authenticated user ID

#### Progress Mapping
- Presentation with step_count: 5
- Verify conceptual mapping: [0.0, 0.25, 0.5, 0.75, 1.0]
- Note: Actual calculation in application layer, not database

## Rollback Plan

If issues discovered:
1. Comment out createPresentationsCollection call
2. Manually drop presentations collection if created
3. Restart application
4. Debug and fix issues
5. Re-enable migration

## Success Criteria

✓ Application starts without errors
✓ presentations collection exists with all fields
✓ Can create presentation via authenticated request
✓ Can list/view presentations publicly
✓ Owner-only updates work correctly
✓ Relations to sync_sessions functional
✓ JSON field accepts array data
✓ Migration is idempotent

## Dependencies

- sync_sessions collection must exist (T-006-01 completed)
- PocketBase v0.36.5 (verified in go.mod)
- No external dependencies

## Risk Mitigation

**Risk:** JSON field validation
**Mitigation:** Document that application layer must validate step_labels array length

**Risk:** Users collection missing
**Mitigation:** Conditional field addition, no ownership tracking if missing

**Risk:** Migration fails mid-execution
**Mitigation:** Transaction handled by PocketBase, automatic rollback

## Notes

- No custom routes needed for basic CRUD (PocketBase handles)
- Active session management likely needs custom routes (future ticket)
- Progress calculation remains application concern, not database
- Consider adding indexes if performance issues arise later