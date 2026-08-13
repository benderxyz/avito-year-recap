import { TextInput, type TextInputProps } from '@mantine/core';
import { Controller, type FieldPath, type FieldValues } from 'react-hook-form';
import type { FormFieldProps } from '../types';

type FormTextInputProps<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = FormFieldProps<TFieldValues, TFieldName> &
  Omit<TextInputProps, 'name' | 'value' | 'defaultValue' | 'onChange' | 'error'>;

export default function FormTextInput<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ name, control, ...props }: FormTextInputProps<TFieldValues, TFieldName>) {
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <TextInput
          {...props}
          {...field}
          value={field.value ?? ''}
          error={fieldState.error?.message}
        />
      )}
    />
  );
}
