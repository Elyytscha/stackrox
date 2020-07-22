import { gql } from '@apollo/client';

export const SUBJECT_WITH_CLUSTER_FRAGMENT = gql`
    fragment subjectWithClusterFields on SubjectWithClusterID {
        id: name
        subject {
            name
            kind
            namespace
        }
        type
        scopedPermissions {
            scope
            permissions {
                key
                values
            }
        }
        clusterAdmin
        roles {
            id
            name
        }
    }
`;

export const SUBJECTS_QUERY = gql`
    query subjects($query: String) {
        clusters {
            id
            name
            subjects(query: $query) {
                ...subjectWithClusterFields
            }
        }
    }
    fragment subjectWithClusterFields on SubjectWithClusterID {
        id: name
        subject {
            name
            kind
            namespace
        }
        type
        clusterAdmin
        roles {
            id
            name
        }
    }
`;

export const SUBJECT_NAME = gql`
    query getSubjectName($clustersQuery: String, $name: String!) {
        clusters(query: $clustersQuery) {
            id
            subject(name: $name) {
                id: name
                subject {
                    name
                }
            }
        }
    }
`;

export const SUBJECT_QUERY = gql`
    query subject($id: String!) {
        clusters {
            id
            name
            subject(name: $id) {
                id: name
                subject {
                    name
                    kind
                    namespace
                }
                type
                scopedPermissions {
                    scope
                    permissions {
                        key
                        values
                    }
                }
                clusterAdmin
                roles {
                    id
                    name
                }
            }
        }
    }
`;
