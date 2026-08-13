import { Textarea, type TextareaProps } from '@mantine/core';
import { Controller, type FieldPath, type FieldValues } from 'react-hook-form';
import type { FormFieldProps } from '../types';

type FormTextareaProps<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = FormFieldProps<TFieldValues, TFieldName> &
  Omit<TextareaProps, 'name' | 'value' | 'defaultValue' | 'onChange' | 'error'>;

export default function FormTextarea<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ name, control, ...props }: FormTextareaProps<TFieldValues, TFieldName>) {
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Textarea
          {...props}
          {...field}
          value={field.value ?? ''}
          error={fieldState.error?.message}
        />
      )}
    />
  );
}
