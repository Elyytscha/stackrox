import { gql } from '@apollo/client';

const GET_PROCESS_COMMENTS_TAGS_COUNT = gql`
    query processCommentsAndTagsCount($key: ProcessNoteKey!) {
        processCommentsCount(key: $key)
        processTagsCount(key: $key)
    }
`;

export default GET_PROCESS_COMMENTS_TAGS_COUNT;
