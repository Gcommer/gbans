import { useEffect, useState } from 'react';
import { apiGetBansGroups, BanGroupQueryFilter, GroupBanRecord } from '../api';
import { logErr } from '../util/errors';

export const useBansGroup = (opts: BanGroupQueryFilter) => {
    const [loading, setLoading] = useState(false);
    const [count, setCount] = useState<number>(0);
    const [data, setData] = useState<GroupBanRecord[]>([]);
    const [error, setError] = useState();

    useEffect(() => {
        const abortController = new AbortController();
        setLoading(true);
        apiGetBansGroups(
            {
                limit: opts.limit,
                offset: opts.offset,
                order_by: opts.order_by,
                desc: opts.desc,
                source_id: opts.source_id,
                target_id: opts.target_id,
                appeal_state: opts.appeal_state,
                deleted: opts.deleted
            },
            abortController
        )
            .then((bans) => {
                setData(bans.data);
                setCount(bans.count);
            })
            .catch((reason) => {
                logErr(reason);
                setError(error);
            })
            .finally(() => {
                setLoading(false);
            });

        return () => abortController.abort();
    }, [
        opts.limit,
        opts.offset,
        opts.order_by,
        opts.desc,
        opts.source_id,
        opts.target_id,
        opts.appeal_state,
        opts.deleted,
        error
    ]);

    return { data, count, loading, error };
};