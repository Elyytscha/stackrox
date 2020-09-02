import { cloneDeep, isNil, isPlainObject, isObject, isEmpty, pickBy, mapValues } from 'lodash';

/**
 * Checks whether the value is empty (null, undefined, empty string, empty array, empty plain object).
 */
const isNilOrEmpty = (v: unknown): boolean => isNil(v) || v === '' || (isObject(v) && isEmpty(v));

/**
 * Removes empty fields from the object traversing deep into fields with object values.
 *
 * @param {object} object any plain object, it'll not be mutated.
 * @param {EmptyValuePredicate} [predicate=isNilOrEmpty] either a given field value is empty
 * @returns {object} returns a deep copy of the original object with empty fields removed
 */
export default function removeEmptyFieldsDeep(obj: object): object {
    const cloned = cloneDeep(obj);
    // deep clean all the fields with values being objects themselves
    const onlyCleanNestedObjects = mapValues(pickBy(cloned, isPlainObject), removeEmptyFieldsDeep);
    // return back fields with non-object values
    const allFields = {
        ...onlyCleanNestedObjects,
        ...pickBy(cloned, (v) => !isPlainObject(v)),
    };
    // filter out empty fields
    return pickBy(allFields, (v) => !isNilOrEmpty(v));
}
