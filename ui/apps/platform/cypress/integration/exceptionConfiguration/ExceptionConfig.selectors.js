export const vulnerabilitiesConfigSelectors = {
    saveButton: 'button:contains("Save")',

    dayOptionInput: (index) => `input[id="expiryOptions.dayOptions[${index}].numDays"]`,
    dayOptionEnabledSwitch: (index) => `input[id="expiryOptions.dayOptions[${index}].enabled"]`,

    indefiniteOptionEnabledSwitch: 'input[id="TODO.enabled"]',
    whenAllCveFixableSwitch: 'input[id="expiryOptions.fixableCveOptions.allFixable"]',
    whenAnyCveFixableSwitch: 'input[id="expiryOptions.fixableCveOptions.anyFixable"]',
    customDateSwitch: 'input[id="expiryOptions.customDate"]',
};
