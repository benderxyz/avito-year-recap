import { Switch, type SwitchProps } from '@mantine/core';
import { Controller, type FieldPath, type FieldValues } from 'react-hook-form';
import type { FormFieldProps } from '../types';

type FormSwitchProps<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> = FormFieldProps<TFieldValues, TFieldName> &
  Omit<
    SwitchProps,
    'name' | 'value' | 'defaultValue' | 'checked' | 'defaultChecked' | 'onChange' | 'error'
  >;

export default function FormSwitch<
  TFieldValues extends FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ name, control, ...props }: FormSwitchProps<TFieldValues, TFieldName>) {
  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <Switch
          {...props}
          name={field.name}
          ref={field.ref}
          checked={Boolean(field.value)}
          onBlur={field.onBlur}
          onChange={(event) => field.onChange(event.currentTarget.checked)}
          error={fieldState.error?.message}
        />
      )}
    />
  );
}
