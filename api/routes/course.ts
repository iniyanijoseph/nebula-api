import express from 'express';

import { courseSearch, courseById } from '../controllers/course';

const router = express.Router();

router.get('/search?', courseSearch);
router.get('/:id', courseById);

export default router;