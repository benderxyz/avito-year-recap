import { z } from 'zod';
import { GetApiAdminPreviewMode } from '@/shared/api/generated/model/getApiAdminPreviewMode';
import { PREVIEW_PALETTE_IDS } from '../lib/palettes';

export const previewFormSchema = z.object({
  reactVersion: z.string().min(1),
  year: z.number().int(),
  mode: z.enum([GetApiAdminPreviewMode.private, GetApiAdminPreviewMode.public]),
  seed: z.number().nullable(),
  themeId: z.enum(PREVIEW_PALETTE_IDS),
  autoplay: z.boolean(),
  loop: z.boolean(),
  gestures: z.boolean(),
  tapNav: z.boolean(),
  holdToPause: z.boolean(),
  reducedMotion: z.boolean(),
});

export type PreviewFormValues = z.infer<typeof previewFormSchema>;

export const previewFormDefaults: PreviewFormValues = {
  reactVersion: '',
  year: new Date().getUTCFullYear(),
  mode: GetApiAdminPreviewMode.private,
  seed: null,
  themeId: 'avitoLight',
  autoplay: false,
  loop: false,
  gestures: true,
  tapNav: false,
  holdToPause: false,
  reducedMotion: false,
};
