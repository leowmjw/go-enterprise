# Batch Processing System Learnings

## Project Overview

This document summarizes the implementation of a batch processing system for payment collection, focusing on robust error handling, exit code management, and proper Bash scripting practices.

## Repository Structure

```
/Users/leow/GOMOD/go-enterprise/
└── internal/
    └── batch/
        ├── scripts/
        │   ├── func_collect_payment.sh     # Function library for payment processing
        │   ├── single_payment_collection.sh # Processes a single payment with error handling
        │   ├── traditional_payment_collection.sh # Batch processes multiple payments
        │   └── http_batch_test.sh          # HTTP-based batch testing script
        ├── bash-expert/                    # Bash best practices reference
        │   └── bash_commands.yaml          # Bash commands reference
        ├── postgres-jit.go                # Go code for database operations
        └── LEARNINGS.md                   # This document
```

## Script Components

### 1. Function Library (`func_collect_payment.sh`)

#### Key Functions

- **`myip()`**: Retrieves the public IP address using curl and processes local IP addresses
- **`process_step1(OrderID)`**: 
  - Simulates a payment processing step with 50% failure rate
  - Returns exit code 1 on failure
  - On success: sleeps 1-3 seconds and outputs "Step1 {OrderID}"

- **`process_step2(OrderID)`**: 
  - More complex failure modes:
    - 20% chance of timeout (exit code 2)
    - 20% chance of gibberish output (exit code 3)
    - 60% chance of success
  - On success: sleeps 0-2 seconds and outputs "Step2 {OrderID}"

#### Implementation Details

```bash
#!/bin/bash
# Function 1: Fails 50% of the time
process_step1() {
    local order_id="$1"
    
    # 50% chance of failure
    if (( RANDOM % 2 )); then
        echo "FAILED: Processing Step 1 for OrderID $order_id"
        return 1
    else
        # Random sleep between 1-3 seconds
        local sleep_time=$(( ( RANDOM % 3 ) + 1 ))
        sleep "$sleep_time"
        echo "Step1 $order_id"
        return 0
    fi
}
```

### 2. Single Payment Processor (`single_payment_collection.sh`)

#### Key Features

- **Error Handling**:
  - Uses `set -euo pipefail` for strict error detection
  - Implements trap-based cleanup for graceful termination
  - Preserves and reports exit codes from processing functions

- **Input Validation**:
  - Verifies OrderID is provided and is a valid number
  - Provides clear usage instructions on error

- **Processing Flow**:
  - Sources the function library
  - Processes steps sequentially with proper error checking
  - Captures and displays detailed error messages

#### Implementation Highlights

```bash
# Setup error handling with cleanup function
cleanup() {
    local exit_code=$?
    echo "Cleaning up resources..."
    
    if [[ $exit_code -ne 0 ]]; then
        if [[ -n "$LAST_ERROR_MSG" ]]; then
            echo "ERROR: Script terminated with exit code: $exit_code - $LAST_ERROR_MSG" >&2
        else
            echo "ERROR: Script terminated with exit code: $exit_code" >&2
        fi
    fi
    
    exit $exit_code
}

# Register for EXIT signal
trap cleanup EXIT
```

### 3. Batch Processor (`traditional_payment_collection.sh`)

#### Key Features

- **Batch Processing**:
  - Processes a predefined list of OrderIDs
  - Continues processing despite individual failures

- **Enhanced Output**:
  - Color-coded output (green for success, red for errors)
  - Timing information for each process
  - Summary statistics (success/failure rates)

- **Error Handling**:
  - Captures and displays all errors without interrupting the batch
  - Preserves original error messages from the single payment processor

#### Implementation Highlights

```bash
# Process each OrderID in the list
for order_id in "${ORDER_IDS[@]}"; do
    # Record start time
    start_time=$(date +%s)
    
    # Call the single payment collection script and capture output
    # We use set +e to prevent the loop from exiting if the script fails
    set +e
    output=$($SOURCE_DIR/single_payment_collection.sh "$order_id" 2>&1)
    exit_code=$?
    set -e
    
    # Record end time and calculate duration
    end_time=$(date +%s)
    duration=$((end_time - start_time))
    
    # Display result based on exit code
    if [[ $exit_code -eq 0 ]]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        echo -e "${GREEN}SUCCESS: OrderID $order_id processed successfully in ${duration}s${NC}"
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
        echo -e "${RED}ERROR: OrderID $order_id failed with exit code $exit_code in ${duration}s${NC}"
    fi
done
```

## Best Practices Implemented

### 1. Error Handling

- **Strict Mode**: Using `set -euo pipefail` to catch errors early
- **Exit Codes**: Proper use and propagation of meaningful exit codes
- **Signal Trapping**: Using `trap` to ensure cleanup on script termination
- **Error Messages**: Detailed error messages sent to stderr

### 2. Variable Handling

- **Quoting**: Proper quoting of variables to prevent word splitting
- **Local Variables**: Using `local` for function-scoped variables
- **Parameter Validation**: Checking parameters before use

### 3. Script Structure

- **Modular Design**: Separating functionality into reusable functions
- **Library Sourcing**: Properly sourcing function libraries
- **Comments**: Clear comments explaining code functionality
- **Consistent Formatting**: Consistent indentation and code style

### 4. Output Formatting

- **Color Coding**: Using ANSI color codes for better readability
- **Structured Output**: Consistent output format with clear section markers
- **Progress Indicators**: Showing progress through batch operations

## Testing Results

The batch processing system was tested with 10 different OrderIDs, demonstrating:

- Approximately 30-40% success rate (matching the designed failure rates)
- Proper error handling and reporting
- Accurate timing information
- Consistent progress through the batch despite failures

## Integration Points

- The scripts can be integrated with Go code in `postgres-jit.go` for database operations
- The batch processor can be extended to read OrderIDs from external sources (files, databases)
- The system can be integrated with monitoring tools via the structured output

## Future Enhancements

1. Add logging to external files
2. Implement parallel processing for better performance
3. Add retry mechanisms for failed processes
4. Integrate with notification systems for critical failures

## Demo Scripts Improvements (March 2025)

A comprehensive review and enhancement of the demo scripts was performed, implementing additional best practices and advanced Bash techniques.

### 1. Key Scripts Improved

- **`demo-start.sh`**: Script to start the SuperScript application
- **`demo-stop.sh`**: Script to safely stop the SuperScript application
- **`demo-1-setup.sh`**: Setup script for the demo environment
- **`demo-2-traditional.sh`**: Demo for traditional payment processing
- **`demo-3-single-payment.sh`**: Demo for single payment workflow
- **`demo-4-orchestrator.sh`**: Demo for orchestrator-based workflow

### 2. Terminal Handling Improvements

- **Portable Color Implementation**:
  ```bash
  # Check if stdout is a terminal before using colors
  if [[ -t 1 ]]; then
      readonly BOLD=$(tput bold)
      readonly GREEN=$(tput setaf 2)
      readonly YELLOW=$(tput setaf 3)
      readonly RED=$(tput setaf 1)
      readonly RESET=$(tput sgr0)
  else
      readonly BOLD=""
      readonly GREEN=""
      readonly YELLOW=""
      readonly RED=""
      readonly RESET=""
  fi
  ```
  
- **Progress Indicators**:
  ```bash
  # Function to display a spinner during process waits
  display_spinner() {
      local pid=$1
      local message="$2"
      local delay=0.1
      local spinstr='|/-\\'
      
      while ps -p "$pid" > /dev/null; do
          local temp=${spinstr#?}
          printf "\r%s %c" "$message" "${spinstr:0:1}"
          spinstr="$temp${spinstr:0:1}"
          sleep "$delay"
      done
      printf "\r%s Done.\n" "$message"
  }
  ```

### 3. Advanced Process Management

- **Graceful Termination with Fallbacks**:
  ```bash
  # Try graceful termination with SIGTERM, fallback to SIGKILL
  kill -15 "${pid}" 2>/dev/null || true
  
  # Wait for termination with timeout
  local counter=0
  while is_app_running && ((counter < STOP_TIMEOUT)); do
      printf "."
      sleep 1
      ((counter++))
  done
  
  # Force kill if still running
  if is_app_running; then
      kill -9 "${pid_updated}" 2>/dev/null || true
  fi
  ```

- **Port and Process Verification**:
  ```bash
  # Multiple fallback methods for checking port usage
  is_port_in_use() {
      if command_exists nc; then
          nc -z localhost "${PORT}" > /dev/null 2>&1
          return $?
      elif command_exists lsof; then
          lsof -i :"${PORT}" > /dev/null 2>&1
          return $?
      else
          # Pure Bash fallback method
          (</dev/tcp/localhost/${PORT}) > /dev/null 2>&1
          return $?
      fi
  }
  ```

### 4. Enhanced User Interaction

- **Robust User Confirmation**:
  ```bash
  ask_for_confirmation() {
      local prompt="$1"
      local default="${2:-n}" # Default is 'n' unless specified
      
      local yn
      while true; do
          printf "%s %s: " "${prompt}" "($(if [[ "${default}" = "y" ]]; then echo "Y/n"; else echo "y/N"; fi))"
          read -r -n 1 yn </dev/tty
          
          case "$yn" in
              "" ) # Enter pressed - use default
                  if [[ "${default}" = "y" ]]; then
                      return 0
                  else
                      return 1
                  fi
                  ;;
              [Yy]* ) return 0 ;;
              [Nn]* ) return 1 ;;
              * ) print_warning "Please answer yes or no." ;;
          esac
      done
  }
  ```

### 5. Script Organization Improvements

- **Script Configuration Section**:
  ```bash
  # Get script directory (safer approach for sourced scripts)
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  
  # Define application configuration
  readonly APP_PROCESS_PATTERN="bin/superscript"
  readonly APP_PORT=8080
  readonly STOP_TIMEOUT=10    # seconds to wait for graceful termination
  readonly FORCE_KILL_DELAY=3 # seconds to wait before force kill
  ```

- **Command Availability Checking**:
  ```bash
  command_exists() {
      command -v "$1" >/dev/null 2>&1
  }
  
  # Check for required commands
  if ! command_exists pgrep; then
      print_error "Required command 'pgrep' not found"
      exit 1
  fi
  ```

### 6. Error Handling Extensions

- **Enhanced Trap Handling**:
  ```bash
  # More comprehensive trap for multiple signals
  trap cleanup EXIT HUP INT TERM
  
  # Preserve original exit code in cleanup
  cleanup() {
      local exit_code=$?
      # Cleanup actions here
      return ${exit_code}  # Preserve original exit code
  }
  ```

- **Improved Error Reporting**:
  ```bash
  print_error() {
      printf "%s%s%s\n" "${RED}${BOLD}" "ERROR: $1" "${RESET}" >&2
  }
  
  print_warning() {
      printf "%s%s%s\n" "${YELLOW}" "WARNING: $1" "${RESET}" >&2
  }
  ```

### 7. Performance Metrics

- **Timing Operations**:
  ```bash
  # Start timing
  local start_time=$(date +%s)
  
  # Perform operations
  stop_app
  
  # Calculate elapsed time
  local end_time=$(date +%s)
  local elapsed=$((end_time - start_time))
  print_success "Operation completed in ${elapsed} seconds!"
  ```

These improvements substantially enhanced the robustness, reliability, and user experience of the demo scripts while maintaining compatibility with different environments.
5. Enhance reporting with more detailed statistics

---

*Last Updated: March 22, 2025*
