import { Select, type SelectProps } from '@mantine/core';

type TriState = boolean | undefined;

type TriStateFilterProps = Omit<SelectProps, 'value' | 'onChange' | 'data'> & {
  value: TriState;
  onChange: (value: TriState) => void;
  trueLabel?: string;
  falseLabel?: string;
  allLabel?: string;
};

function toSelectValue(value: TriState) {
  if (value === true) {
    return 'true';
  }
  if (value === false) {
    return 'false';
  }
  return null;
}

function fromSelectValue(value: string | null): TriState {
  if (value === 'true') {
    return true;
  }
  if (value === 'false') {
    return false;
  }
  return undefined;
}

export default function TriStateFilter({
  value,
  onChange,
  label = 'Enabled',
  trueLabel = 'Yes',
  falseLabel = 'No',
  allLabel = 'All',
  ...props
}: TriStateFilterProps) {
  return (
    <Select
      {...props}
      label={label}
      clearable
      placeholder={allLabel}
      value={toSelectValue(value)}
      onChange={(next) => onChange(fromSelectValue(next))}
      data={[
        { value: 'true', label: trueLabel },
        { value: 'false', label: falseLabel },
      ]}
    />
  );
}
