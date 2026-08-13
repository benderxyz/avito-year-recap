import { Select, type SelectProps } from '@mantine/core';
import { Controller, type FieldPath, type FieldValues } from 'react-hook-form';
import type { FormFieldProps } from '../types';

type FormSelectProps<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = FormFieldProps<TFieldValues, TFieldName> &
  Omit<SelectProps, 'name' | 'value' | 'defaultValue' | 'onChange' | 'error'>;

export default function FormSelect<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ name, control, ...props }: FormSelectProps<TFieldValues, TFieldName>) {
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Select
          {...props}
          name={field.name}
          ref={field.ref}
          value={field.value ?? null}
          onBlur={field.onBlur}
          onChange={field.onChange}
          error={fieldState.error?.message}
        />
      )}
    />
  );
}
