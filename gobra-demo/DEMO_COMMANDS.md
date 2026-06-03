# Gobra Plugin gobra-demo - Command List

## Setup Commands

```bash
# Generate ANTLR parser and build framework
cd framework && make build
```

## Success Examples - Contracts

```bash
# Verify all functions with correct contracts
cd framework && ./build/framework file --plugin gobra --plugin-config gobra:backend:docker ../gobra-demo/success/contracts.go
```

**Expected Result:** `Verification succeed`

## Success Examples - Assertions

```bash
# Verify all functions with correct assertions
cd framework && ./build/framework file --plugin gobra --plugin-config gobra:backend:docker ../gobra-demo/success/asserts.go
```

**Expected Result:** `Verification succeed`

## Failure Examples - Assertions

```bash
# Verify functions with incorrect assertions (should fail)
cd framework && ./build/framework file --plugin gobra --plugin-config gobra:backend:docker ../gobra-demo/failure/asserts.go
```

**Expected Result:**
```
Verification failed
	Error: condition at 10:13 may fail
		x == x + 1
	Error: condition at 4:13 may fail
		x + y >= x
```