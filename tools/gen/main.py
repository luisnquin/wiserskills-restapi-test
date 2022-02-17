from random import choice, randint
from faker import Faker
from sys import argv


fake: object = Faker()

participant_ids: list[int] = []
event_ids: list[int] = []

used_pairs: list[tuple[int]] = []


def gen_event() -> str:
    stmt: str = 'INSERT INTO events(name) VALUES(\'%s\');\n'

    event_ids.append(1) if \
        len(event_ids) == 0 else event_ids.append(event_ids[-1]+1)

    return stmt % (fake.unique.name())


def gen_participant() -> str:
    stmt: str = 'INSERT INTO participants(firstname, lastname, age) VALUES(\'%s\', \'%s\', %d);\n'
    name: list[str] = fake.unique.name().split(' ')

    participant_ids.append(1) if \
        len(participant_ids) == 0 else participant_ids.append(participant_ids[-1]+1)

    return stmt % (name[0], name[1], randint(18, 129))


def gen_ticket() -> str:
    stmt: str = 'INSERT INTO tickets(event, participant) VALUES(%d, %d);\n'
    pairs: tuple[int] = (choice(event_ids), choice(participant_ids))
    if pairs in used_pairs:
        gen_ticket()

    used_pairs.append(pairs)
    return stmt % pairs


def main(loops: int):
    events: tuple[str] = tuple(gen_event() for _ in range(loops))
    participants: tuple[str] = tuple(gen_participant() for _ in range(loops*5))
    tickets: tuple[str] = tuple(gen_ticket() for _ in range(loops*8))

    with open(file='./../out/statements.sql', mode='w', encoding='UTF-8') as file:
        for group in (events, participants, tickets):
            file.writelines(group)


if __name__ == '__main__':
    main(loops=10 if len(argv) == 1 else int(argv[1]))
