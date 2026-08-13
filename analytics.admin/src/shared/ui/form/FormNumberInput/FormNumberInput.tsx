import { NumberInput, type NumberInputProps } from '@mantine/core';
import { Controller, type FieldPath, type FieldValues } from 'react-hook-form';
import type { FormFieldProps } from '../types';

type FormNumberInputProps<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = FormFieldProps<TFieldValues, TFieldName> &
  Omit<NumberInputProps, 'name' | 'value' | 'defaultValue' | 'onChange' | 'error'>;

export default function FormNumberInput<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ name, control, ...props }: FormNumberInputProps<TFieldValues, TFieldName>) {
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <NumberInput
          {...props}
          name={field.name}
          ref={field.ref}
          value={field.value ?? ''}
          onBlur={field.onBlur}
          onChange={field.onChange}
          error={fieldState.error?.message}
        />
      )}
    />
  );
}
