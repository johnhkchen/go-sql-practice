/**
 * Step-Progress Conversion Utilities
 *
 * Provides conversion between discrete step indices (0-based) and continuous
 * progress values (0.0-1.0) for presentation navigation. Based on formulas
 * from backend presentations.go:95
 */

export interface StepNavigationState {
  currentStep: number;
  canGoPrevious: boolean;
  canGoNext: boolean;
  totalSteps: number;
}

/**
 * Convert a step index to a progress value (0.0-1.0)
 * @param stepIndex 0-based step index
 * @param stepCount Total number of steps in presentation
 * @returns Progress value between 0.0 and 1.0
 */
export function stepToProgress(stepIndex: number, stepCount: number): number {
  if (stepCount <= 1) {
    return 0.0;
  }

  // Ensure stepIndex is within valid bounds
  const clampedIndex = Math.max(0, Math.min(stepIndex, stepCount - 1));

  return clampedIndex / (stepCount - 1);
}

/**
 * Convert a progress value to a step index
 * @param progress Progress value between 0.0 and 1.0
 * @param stepCount Total number of steps in presentation
 * @returns 0-based step index
 */
export function progressToStep(progress: number, stepCount: number): number {
  if (stepCount <= 1) {
    return 0;
  }

  // Ensure progress is within valid bounds
  const clampedProgress = Math.max(0.0, Math.min(progress, 1.0));

  return Math.round(clampedProgress * (stepCount - 1));
}

/**
 * Validate if a step index is within valid bounds
 * @param stepIndex Step index to validate
 * @param stepCount Total number of steps
 * @returns True if step index is valid
 */
export function validateStepIndex(stepIndex: number, stepCount: number): boolean {
  return Number.isInteger(stepIndex) &&
         stepIndex >= 0 &&
         stepIndex < stepCount;
}

/**
 * Get navigation state for a given step
 * @param currentStep Current step index
 * @param stepCount Total number of steps
 * @returns Navigation state with boundary information
 */
export function getNavigationState(currentStep: number, stepCount: number): StepNavigationState {
  return {
    currentStep,
    canGoPrevious: currentStep > 0,
    canGoNext: currentStep < stepCount - 1,
    totalSteps: stepCount,
  };
}

/**
 * Get the next step index, respecting boundaries
 * @param currentStep Current step index
 * @param stepCount Total number of steps
 * @returns Next step index or current if at boundary
 */
export function getNextStep(currentStep: number, stepCount: number): number {
  return Math.min(currentStep + 1, stepCount - 1);
}

/**
 * Get the previous step index, respecting boundaries
 * @param currentStep Current step index
 * @param stepCount Total number of steps
 * @returns Previous step index or current if at boundary
 */
export function getPreviousStep(currentStep: number, stepCount: number): number {
  return Math.max(currentStep - 1, 0);
}

/**
 * Format step display text with optional labels
 * @param stepIndex Current step index (0-based)
 * @param stepCount Total number of steps
 * @param stepLabels Optional array of step labels
 * @returns Formatted step text like "Step 3 of 5 — Label"
 */
export function formatStepDisplay(
  stepIndex: number,
  stepCount: number,
  stepLabels?: string[]
): string {
  const stepNumber = stepIndex + 1; // Convert to 1-based for display
  let display = `Step ${stepNumber} of ${stepCount}`;

  if (stepLabels && stepLabels[stepIndex]) {
    display += ` — ${stepLabels[stepIndex]}`;
  }

  return display;
}