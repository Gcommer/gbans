import NiceModal, { muiDialogV5, useModal } from '@ebay/nice-modal-react';
import PersonIcon from '@mui/icons-material/Person';
import { Dialog, DialogActions, DialogContent, DialogTitle } from '@mui/material';
import Grid from '@mui/material/Unstable_Grid2';
import { useForm } from '@tanstack/react-form';
import { useMutation } from '@tanstack/react-query';
import { z } from 'zod';
import { apiNewsCreate, apiNewsSave, NewsEntry } from '../../api/news.ts';
import { useUserFlashCtx } from '../../hooks/useUserFlashCtx.ts';
import { Heading } from '../Heading';
import { Buttons } from '../field/Buttons.tsx';
import { CheckboxSimple } from '../field/CheckboxSimple.tsx';
import { MarkdownField } from '../field/MarkdownField.tsx';
import { TextFieldSimple } from '../field/TextFieldSimple.tsx';

export const NewsEditModal = NiceModal.create(({ entry }: { entry?: NewsEntry }) => {
    const modal = useModal();
    const { sendError, sendFlash } = useUserFlashCtx();

    const mutation = useMutation({
        mutationKey: ['newsEdit'],
        mutationFn: async (values: { title: string; body_md: string; is_published: boolean }) => {
            if (entry?.news_id) {
                return await apiNewsSave({ ...entry, ...values });
            } else {
                return await apiNewsCreate(values.title, values.body_md, values.is_published);
            }
        },
        onSuccess: async (entry) => {
            modal.resolve(entry);
            sendFlash('success', 'News edited successfully.');
            await modal.hide();
        },
        onError: sendError
    });

    const { Field, Subscribe, handleSubmit, reset } = useForm({
        onSubmit: async ({ value }) => {
            mutation.mutate(value);
        },
        validators: {
            onChange: z.object({
                title: z.string().min(5, 'Min length 5'),
                body_md: z.string().min(10, 'Min length 10'),
                is_published: z.boolean()
            })
        },
        defaultValues: {
            title: entry?.title ?? '',
            body_md: entry?.body_md ?? '',
            is_published: entry?.is_published ?? false
        }
    });

    return (
        <Dialog {...muiDialogV5(modal)} fullWidth maxWidth={'sm'}>
            <form
                onSubmit={async (e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    await handleSubmit();
                }}
            >
                <DialogTitle component={Heading} iconLeft={<PersonIcon />}>
                    News Editor
                </DialogTitle>
                <DialogContent>
                    <Grid container spacing={2}>
                        <Grid xs={12}>
                            <Field
                                name={'title'}
                                children={(props) => {
                                    return <TextFieldSimple {...props} label={'Title'} />;
                                }}
                            />
                        </Grid>
                        <Grid xs={12}>
                            <Field
                                name={'body_md'}
                                children={(props) => {
                                    return <MarkdownField {...props} label={'Body'} />;
                                }}
                            />
                        </Grid>
                        <Grid xs={12}>
                            <Field
                                name={'is_published'}
                                children={({ state, handleBlur, handleChange }) => {
                                    return (
                                        <CheckboxSimple
                                            checked={state.value}
                                            onChange={(_, v) => handleChange(v)}
                                            onBlur={handleBlur}
                                            label={'Is Published'}
                                        />
                                    );
                                }}
                            />
                        </Grid>
                    </Grid>
                </DialogContent>
                <DialogActions>
                    <Grid container>
                        <Grid xs={12} mdOffset="auto">
                            <Subscribe
                                selector={(state) => [state.canSubmit, state.isSubmitting]}
                                children={([canSubmit, isSubmitting]) => {
                                    return (
                                        <Buttons
                                            reset={reset}
                                            canSubmit={canSubmit}
                                            isSubmitting={isSubmitting}
                                            onClose={async () => {
                                                await modal.hide();
                                            }}
                                        />
                                    );
                                }}
                            />
                        </Grid>
                    </Grid>
                </DialogActions>
            </form>
        </Dialog>
    );
});
