import React from 'react';
import PropTypes from 'prop-types';
import { Field } from 'redux-form';
import NumericInput from 'react-numeric-input';

const ReduxNumericInput = (props) => (
    <NumericInput
        max={props.max}
        min={props.min}
        step={props.step}
        key={props.input.value}
        field={props.input.value}
        id={props.input.value}
        value={props.input.value}
        placeholder={props.placeholder}
        onBlur={props.input.onChange}
        noStyle
        className={`${props.disabled ? 'bg-base-200' : 'hover:border-base-400'} ${props.className}`}
        disabled={props.disabled}
    />
);

ReduxNumericInput.propTypes = {
    input: PropTypes.shape({
        value: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
        onChange: PropTypes.func,
    }).isRequired,
    placeholder: PropTypes.string.isRequired,
    min: PropTypes.number.isRequired,
    max: PropTypes.number.isRequired,
    step: PropTypes.number.isRequired,
    className: PropTypes.string.isRequired,
    disabled: PropTypes.bool.isRequired,
};

const ReduxNumericInputField = ({ id, name, min, max, placeholder, step, className, disabled }) => (
    <Field
        id={id}
        key={name}
        name={name}
        placeholder={placeholder}
        parse={parseFloat}
        min={min}
        max={max}
        step={step}
        component={ReduxNumericInput}
        className={className}
        disabled={disabled}
    />
);

ReduxNumericInputField.propTypes = {
    id: PropTypes.string,
    name: PropTypes.string.isRequired,
    min: PropTypes.number,
    max: PropTypes.number,
    step: PropTypes.number,
    placeholder: PropTypes.string,
    className: PropTypes.string,
    disabled: PropTypes.bool,
};

ReduxNumericInputField.defaultProps = {
    id: null,
    min: 1,
    max: Number.MAX_SAFE_INTEGER,
    step: 1,
    placeholder: '',
    className: 'bg-base-100 border-2 rounded-l p-3 text-base-600 border-base-300 w-full font-600',
    disabled: false,
};

export default ReduxNumericInputField;
