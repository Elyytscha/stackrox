import React from 'react';
import {
    PageSection,
    Pagination,
    Toolbar,
    ToolbarContent,
    ToolbarItem,
} from '@patternfly/react-core';
import { TableComposable, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import { Link } from 'react-router-dom';

import { exceptionManagementPath } from 'routePaths';
import useURLPagination from 'hooks/useURLPagination';

function PendingApprovals() {
    const { page, perPage, setPage, setPerPage } = useURLPagination(20);

    return (
        <PageSection>
            <Toolbar>
                <ToolbarContent>
                    <ToolbarItem variant="pagination" alignment={{ default: 'alignRight' }}>
                        <Pagination
                            itemCount={1}
                            page={page}
                            perPage={perPage}
                            onSetPage={(_, newPage) => setPage(newPage)}
                            onPerPageSelect={(_, newPerPage) => setPerPage(newPerPage)}
                            isCompact
                        />
                    </ToolbarItem>
                </ToolbarContent>
            </Toolbar>
            <TableComposable borders={false}>
                <Thead noWrap>
                    <Tr>
                        <Th>Request ID</Th>
                        <Th>Requester</Th>
                        <Th>Requested action</Th>
                        <Th>Requested</Th>
                        <Th>Expires</Th>
                        <Th>Scope</Th>
                        <Th>Requested items</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    <Tr>
                        <Td>
                            <Link to={`${exceptionManagementPath}/requests/:MS-D-20230807-233`}>
                                MS-D-20230807-233
                            </Link>
                        </Td>
                        <Td>Mansur Sayed</Td>
                        <Td>Deferral (30 days)</Td>
                        <Td>7/1/2023</Td>
                        <Td>8/7/23</Td>
                        <Td>gcr.io/ultra-current-825/srox/asset-cache:latest-v2</Td>
                        <Td>6 CVEs</Td>
                    </Tr>
                </Tbody>
            </TableComposable>
        </PageSection>
    );
}

export default PendingApprovals;
