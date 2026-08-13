import type { Control, FieldPath, FieldValues } from 'react-hook-form';

export interface FormFieldProps<
  TFieldValues extends FieldValues = FieldValues,
  TFieldName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> {
  name: TFieldName;
  control: Control<TFieldValues>;
}
