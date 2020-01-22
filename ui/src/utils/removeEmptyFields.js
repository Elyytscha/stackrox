import flatten from 'flat';
import omitBy from 'lodash/omitBy';

export default function removeEmptyFields(obj) {
    const flattenedObj = flatten(obj);
    const omittedObj = omitBy(
        flattenedObj,
        value =>
            value === null ||
            value === undefined ||
            value === '' ||
            value === [] ||
            (Array.isArray(value) && !value.length)
    );
    const newObj = flatten.unflatten(omittedObj);

    // The following fields are not used if they have falsy values,
    //   but those still returned from the API,
    //   so we have to filter them out separately
    const exceptionFields = [
        'imageAgeDays',
        'scanAgeDays',
        'noScanExists',
        'readOnlyRootFs',
        'whitelistEnabled'
    ];
    exceptionFields.forEach(fieldName => {
        if (!newObj[fieldName]) {
            delete newObj[fieldName];
        }
    });

    return newObj;
}
