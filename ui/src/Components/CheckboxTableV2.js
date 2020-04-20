import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactTablePropTypes from 'react-table/lib/propTypes';
import TableV2 from 'Components/TableV2';

class CheckboxTable extends Component {
    static propTypes = {
        columns: ReactTablePropTypes.columns.isRequired,
        rows: PropTypes.arrayOf(PropTypes.object).isRequired,
        toggleRow: PropTypes.func.isRequired,
        toggleSelectAll: PropTypes.func.isRequired,
        selection: PropTypes.arrayOf(PropTypes.string),
    };

    static defaultProps = {
        selection: [],
    };

    setTableRef = (table) => {
        this.reactTable = table;
    };

    toggleRowHandler = ({ id }) => () => {
        this.props.toggleRow(id);
    };

    stopPropagationOnClick = (e) => e.stopPropagation();

    toggleSelectAllHandler = () => {
        this.props.toggleSelectAll();
    };

    someSelected = () => {
        const { selection, rows } = this.props;
        return selection.length !== 0 && selection.length < rows.length;
    };

    allSelected = () => {
        const { selection, rows } = this.props;
        return selection.length !== 0 && selection.length === rows.length;
    };

    addCheckboxColumns = () => {
        const { columns, selection } = this.props;
        return [
            {
                id: 'checkbox',
                accessor: '',
                Cell: ({ original }) => (
                    <input
                        type="checkbox"
                        checked={selection.includes(original.id)}
                        onChange={this.toggleRowHandler(original)}
                        onClick={this.stopPropagationOnClick} // don't want checkbox click to select the row
                        aria-label="Toggle row select"
                    />
                ),
                Header: () => (
                    <input
                        type="checkbox"
                        checked={this.allSelected()}
                        ref={(input) => {
                            if (input) {
                                input.indeterminate = this.someSelected(); // eslint-disable-line no-param-reassign
                            }
                        }}
                        onChange={this.toggleSelectAllHandler}
                        aria-label="Toggle select all rows"
                    />
                ),
                sortable: false,
                width: 28,
            },
            ...columns,
        ];
    };

    render() {
        const { ...rest } = this.props;
        const columns = this.addCheckboxColumns();
        return <TableV2 {...rest} columns={columns} setTableRef={this.setTableRef} />;
    }
}

export default CheckboxTable;
