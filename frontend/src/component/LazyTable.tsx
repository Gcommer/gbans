import { defaultRenderer, HeadingCell, Order } from './DataTable';
import React from 'react';
import TableContainer from '@mui/material/TableContainer';
import Table from '@mui/material/Table';
import useTheme from '@mui/material/styles/useTheme';
import TableBody from '@mui/material/TableBody';
import TableRow from '@mui/material/TableRow';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { TableSortLabel } from '@mui/material';

export interface LazyTableProps<T> {
    columns: HeadingCell<T>[];
    sortColumn: keyof T;
    onSortColumnChanged: (column: keyof T) => void;
    onSortOrderChanged: (order: Order) => void;
    sortOrder: Order;
    rows: T[];
}

export interface TableBodyRows<T> {
    columns: HeadingCell<T>[];
    rows: T[];
}

export const LazyTableBody = <T,>({ rows, columns }: TableBodyRows<T>) => {
    return (
        <TableBody>
            {rows.map((row, idx) => {
                return (
                    <TableRow key={`row-${idx}`} hover>
                        {columns.map((col: HeadingCell<T>, colIdx) => {
                            const value = (col?.renderer ?? defaultRenderer)(
                                row,
                                row[col.sortKey as keyof T],
                                col?.sortType ?? 'string'
                            );
                            return (
                                <TableCell
                                    variant="body"
                                    key={`col-${colIdx}`}
                                    align={col?.align ?? 'right'}
                                    onClick={() => {
                                        if (col.onClick) {
                                            col.onClick(row);
                                        }
                                    }}
                                    sx={{
                                        '&:hover': {
                                            cursor: 'pointer'
                                        },
                                        width: col?.width ?? 'auto'
                                    }}
                                >
                                    {value}
                                </TableCell>
                            );
                        })}
                    </TableRow>
                );
            })}
        </TableBody>
    );
};

export interface LazyTableHeaderProps<T> {
    columns: HeadingCell<T>[];
    bgColor: string;
    onSortColumnChanged: (column: keyof T) => void;
    onSortOrderChanged: (order: Order) => void;
    sortColumn: keyof T;
    order: Order;
}

export const LazyTableHeader = <T,>({
    columns,
    bgColor,
    sortColumn,
    order,
    onSortColumnChanged,
    onSortOrderChanged
}: LazyTableHeaderProps<T>) => {
    return (
        <TableHead>
            <TableRow>
                {columns.map((col) => {
                    return (
                        <TableCell
                            variant="head"
                            align={col.align ?? 'right'}
                            key={col.label}
                            sortDirection={order}
                            sx={{
                                width: col?.width ?? 'auto',
                                backgroundColor: bgColor
                            }}
                        >
                            {col.sortable ? (
                                <TableSortLabel
                                    title={col.tooltip}
                                    active={sortColumn === col.sortKey}
                                    direction={order}
                                    onClick={() => {
                                        if (col.sortKey) {
                                            if (sortColumn == col.sortKey) {
                                                onSortOrderChanged(
                                                    order == 'asc'
                                                        ? 'desc'
                                                        : 'asc'
                                                );
                                            } else {
                                                onSortColumnChanged(
                                                    col.sortKey
                                                );
                                            }
                                        }
                                    }}
                                >
                                    <Typography
                                        padding={0}
                                        sx={{
                                            textDecoration:
                                                col.sortKey != sortColumn
                                                    ? 'none'
                                                    : order == 'asc'
                                                    ? 'underline'
                                                    : 'overline'
                                        }}
                                        variant={'button'}
                                    >
                                        {col.label}
                                    </Typography>
                                </TableSortLabel>
                            ) : (
                                <Typography
                                    padding={0}
                                    sx={{
                                        textDecoration:
                                            col.sortKey != sortColumn
                                                ? 'none'
                                                : order == 'asc'
                                                ? 'underline'
                                                : 'overline'
                                    }}
                                    variant={'button'}
                                >
                                    {col.label}
                                </Typography>
                            )}
                        </TableCell>
                    );
                })}
            </TableRow>
        </TableHead>
    );
};

export const LazyTable = <T,>({
    columns,
    sortOrder,
    sortColumn,
    rows,
    onSortColumnChanged,
    onSortOrderChanged
}: LazyTableProps<T>) => {
    const theme = useTheme();
    // eslint-disable-next-line @typescript-eslint/no-unused-vars

    return (
        <TableContainer component={Paper}>
            <Table size={'small'}>
                <LazyTableHeader<T>
                    columns={columns}
                    sortColumn={sortColumn}
                    onSortColumnChanged={onSortColumnChanged}
                    order={sortOrder}
                    bgColor={theme.palette.background.paper}
                    onSortOrderChanged={onSortOrderChanged}
                />
                <LazyTableBody rows={rows} columns={columns} />
            </Table>
        </TableContainer>
    );
};