import { Badge, Card, Stack } from "react-bootstrap";
import { Link } from 'react-router-dom';
import { NoteRenderData } from "../types/Note";

function NoteCard({ id, header, tags }: NoteRenderData) {
    return (
        <Card
            as={Link} to={`/${id}`}
            className="h-100 text-reset text-decoration-none">
            <Card.Body>
                <Stack gap={2} className="align-items-center">
                    <span className="flex-wrap justify-content-center">{header}</span>
                    <Stack gap={1} direction="horizontal" className="mb-2 flex-wrap justify-content-center">
                        {tags.map(tag => (
                            <Badge className="text-truncate" key={tag.id}>
                                {tag.label}
                            </Badge>
                        ))}
                    </Stack>
                </Stack>
            </Card.Body>
        </Card>
    );
}

export default NoteCard